package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vivangkumar/skv/internal/backend"
	"github.com/vivangkumar/skv/internal/wire"
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

// Server represents an skv TCP server.
type Server struct {
	// The backing store for the server.
	backend be

	// listener on which the server listens to
	// for new incoming TCP connections.
	listener net.Listener

	// stop signals for the server to exit.
	stop chan struct{}

	once      sync.Once
	isStopped bool

	// conns tracks the accepted connections.
	conns sync.WaitGroup

	// server specific logger.
	l *log.Entry
}

// New creates a new skv TCP server.
func New(backend be, cfg Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("server: listen: %w", err)
	}

	return &Server{
		backend:   backend,
		listener:  l,
		stop:      make(chan struct{}),
		once:      sync.Once{},
		isStopped: false,
		conns:     sync.WaitGroup{},
		l:         log.WithField("component", "server"),
	}, nil
}

// Listen accepts connections on the server's connection port.
func (s *Server) Listen(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stop:
				return
			default:
				s.l.WithError(err).Error("failed to accept connection")
			}
		}

		conn.SetDeadline(time.Now().Add(idleTimeoutDuration))
		go s.handleConn(ctx, conn)
	}
}

// Stop shuts down the server gracefully.
//
// It does so by closing the TCP listener and
// waiting for all currently served connections to
// be served or to be closed.
//
// It can be safely called multiple times.
func (s *Server) Stop() {
	if s.isStopped {
		return
	}

	s.l.Info("stopping server")

	close(s.stop)
	if err := s.listener.Close(); err != nil {
		s.l.WithError(err).Error("failed to close listener")
	}
	s.conns.Wait()

	s.once.Do(func() {
		s.isStopped = true
	})

	if err := s.backend.Stop(); err != nil {
		s.l.WithError(err).Error("failed to stop backend")
	}
}

// handleConn handles a single TCP connection
// it decodes incoming requests and validates them
// before passing it on to the backend.
func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	s.conns.Add(1)
	defer s.conns.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	s.l.WithField("remote_addr", conn.RemoteAddr().String()).Info("handling connection")

	lr := &io.LimitedReader{R: conn, N: maxReadBytes}
	scanner := bufio.NewScanner(lr)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		cmd, err := wire.DecodeCmd(line)
		if err != nil {
			s.l.WithError(err).Error("failed to decode command")
			s.writeReply(
				conn,
				wire.Err,
				wire.CmdDecodeErr,
				err.Error(),
			)

			continue
		}

		r, err := s.command(ctx, cmd)
		if err != nil {
			s.l.WithError(err).WithField("command", cmd.Type).Error("failed to process command")
			s.writeReply(conn, wire.Err, wire.CmdFailedErr, "command failed")
			continue
		}

		s.writeReply(conn, wire.OK, r.args...)
	}
}

// command processes received commands by forwarding it to the backend.
func (s *Server) command(
	ctx context.Context,
	cmd wire.Cmd,
) (result, error) {
	switch cmd.Type {
	case wire.Set:
		err := s.backend.Set(
			ctx,
			backend.SetCmd{Key: cmd.Args[0], Value: cmd.Args[1]},
		)
		if err != nil {
			return result{}, newCmdError(err)
		}
	case wire.Get:
		v, err := s.backend.Get(ctx, backend.GetCmd{Key: cmd.Args[0]})
		if err != nil {
			return result{}, newCmdError(err)
		}

		return result{args: []string{v}}, nil
	case wire.Del:
		err := s.backend.Delete(ctx, backend.DelCmd{Key: cmd.Args[0]})
		if err != nil {
			return result{}, newCmdError(err)
		}
	}

	return result{}, nil
}

// writeReply writes a reply to the connection.
//
// If a reply cannot be written, then we can assume
// the connection has failed and that the client will
// retry the request again.
func (s *Server) writeReply(
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
