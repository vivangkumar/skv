package node

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/vivangkumar/skv/pkg/backend"
	"github.com/vivangkumar/skv/pkg/wire"
)

const (
	idleTimeoutDuration       = 60 * time.Second
	maxReadBytes        int64 = 1000 // 1KB
)

type be interface {
	Set(ctx context.Context, cmd backend.SetCmd) error
	Get(ctx context.Context, cmd backend.GetCmd) (string, error)
	Delete(ctx context.Context, cmd backend.DelCmd) error
	Stop() error
}

// Node represents an skv cache node.
//
// A node is an abstraction that links
// the storage layer to the network layer
// and is respnsible for translating client
// network requests into command and replying back.
//
// Each node is assigned a unique ID since there can
// be multiple such nodes.
type Node struct {
	// A unique ID for the Node.
	id string

	// The backing store for the node
	backend be

	// listener on which the node listens to
	// for new incoming TCP connections.
	listener net.Listener

	// stop signals for the node to exit.
	stop chan struct{}

	m         sync.Mutex
	isStopped bool

	// conns tracks the accepted connections.
	conns sync.WaitGroup

	// node specific logger
	l *log.Entry
}

// NewNode creates a new instance of Node
func NewNode(backend be, cfg Config) (*Node, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("node: listen: %w", err)
	}

	id := uuid.New().String()

	return &Node{
		backend:   backend,
		listener:  l,
		stop:      make(chan struct{}),
		m:         sync.Mutex{},
		isStopped: false,
		conns:     sync.WaitGroup{},
		id:        id,
		l:         log.WithField("node_id", id),
	}, nil
}

// Listen accepts connections on the node's
// connection port.
func (n *Node) Listen(ctx context.Context) {
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			select {
			case <-n.stop:
				return
			default:
				n.l.WithError(err).Error("node: listener: accept")
			}
		}

		conn.SetDeadline(time.Now().Add(idleTimeoutDuration))
		go n.handleConn(ctx, conn)
	}
}

// ID returns the unique node ID.
func (n *Node) ID() string {
	return n.id
}

// Stop shuts down the node gracefully.
//
// It does so by closing the TCP listener and
// waiting for all currently served connections to
// be served or to be closed.
//
// It can be safely called multiple times.
func (n *Node) Stop() {
	if n.isStopped {
		return
	}

	n.l.WithField("node_id", n.id).Info("node: stopping")

	close(n.stop)
	if err := n.listener.Close(); err != nil {
		n.l.WithError(err).Error("node: listener: close")
	}
	n.conns.Wait()

	n.m.Lock()
	n.isStopped = true
	n.m.Unlock()

	if err := n.backend.Stop(); err != nil {
		n.l.WithError(err).Error("node: backend: stop")
	}
}

// handleConn handles a single TCP connection
// it decodes incoming requests and validates them
// before passing it on to the backend.
func (n *Node) handleConn(ctx context.Context, conn net.Conn) {
	n.conns.Add(1)
	defer n.conns.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	n.l.WithField("remote_addr", conn.RemoteAddr().String()).Info("handling connection")

	lr := &io.LimitedReader{R: conn, N: maxReadBytes}
	scanner := bufio.NewScanner(lr)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		cmd, err := wire.DecodeCmd(line)
		if err != nil {
			n.l.WithError(err).Error("node: handle conn: decode")
			n.writeReply(
				conn,
				wire.Err,
				wire.CmdDecodeErr,
				err.Error(),
			)

			continue
		}

		r, err := n.command(ctx, cmd)
		if err != nil {
			n.l.WithError(err).Error("node: handle conn: command")
			n.writeReply(conn, wire.Err, wire.CmdFailedErr, "command failed")
		}

		n.writeReply(conn, wire.OK, r.args...)
	}
}

type result struct {
	args []string
}

type cmdError struct {
	err error
}

func newCmdError(err error) cmdError {
	return cmdError{err: fmt.Errorf("command: %s", err)}
}

func (c cmdError) CommandFailed() bool {
	return true
}

func (c cmdError) Error() string {
	return c.err.Error()
}

func (n *Node) command(
	ctx context.Context,
	cmd wire.Cmd,
) (result, error) {
	switch cmd.Type {
	case wire.Set:
		err := n.backend.Set(
			ctx,
			backend.SetCmd{Key: cmd.Args[0], Value: cmd.Args[1]},
		)
		if err != nil {
			return result{}, newCmdError(err)
		}
	case wire.Get:
		v, err := n.backend.Get(ctx, backend.GetCmd{Key: cmd.Args[0]})
		if err != nil {
			return result{}, newCmdError(err)
		}

		return result{args: []string{v}}, nil
	case wire.Del:
		err := n.backend.Delete(ctx, backend.DelCmd{Key: cmd.Args[0]})
		if err != nil {
			return result{}, newCmdError(err)
		}
	}

	return result{}, nil
}

// writeReply writes a reply to the connection.
//
// if a reply cannot be written, then we can assume
// the connection has failed and that the client will
// retry the request again.
//
// TODO: metrics.
func (n *Node) writeReply(
	conn net.Conn,
	typ wire.Reply,
	args ...string,
) error {
	var enc []byte

	r := wire.CmdReply{Type: typ, Args: args}
	enc, err := wire.EncodeReply(r)
	if err != nil {
		return fmt.Errorf("write reply: encode: %w", err)
	}

	_, err = conn.Write(enc)
	if err != nil {
		return fmt.Errorf("write reply: write conn: %w", err)
	}

	// Renew the connection deadline after the latest write.
	conn.SetDeadline(time.Now().Add(idleTimeoutDuration))

	return nil
}
