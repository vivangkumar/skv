package wire

import (
	"errors"
	"strings"
)

const (
	errReply = "err"
	okReply  = "ok"

	sep  = ":"
	crlf = "\r\n"
)

// EncodeReply translates CmdReply to a byte array
// so that it can be written over the wire.
func EncodeReply(reply CmdReply) ([]byte, error) {
	switch reply.Type {
	case OK:
		return ok(reply.Args...), nil
	case Err:
		return err(reply.Args...), nil
	default:
		return nil, errors.New("unknown reply type")
	}
}

func ok(args ...string) []byte {
	return reply(okReply, encodeNulls(args...)...)
}

func err(args ...string) []byte {
	return reply(errReply, args...)
}

func encodeNulls(args ...string) []string {
	var withNulls []string
	for _, a := range args {
		if a == "" {
			a = "null"
		}

		withNulls = append(withNulls, a)
	}

	return withNulls
}

func reply(reply string, args ...string) []byte {
	arr := []string{reply}
	if len(args) > 0 {
		arr = append(arr, args...)
	}

	return []byte(terminate(strings.Join(arr, sep)))
}

// apply a terminating \r\n to the end of our replies
func terminate(reply string) string {
	return reply + crlf
}
