package wire

import (
	"errors"
	"fmt"
	"strings"
)

// The message protocol is fairly simple and dumb.
//
// "<cmd>:<arg1>:<arg2>\r\n"
//
// All commands are terminated with \r\n.
//
// For the time being pipelining is not supported
// everything except the first command will be
// discarded.

// DecodeCmd decodes requests sent over the wire
// into internal wire datatypes.
func DecodeCmd(req string) (Cmd, error) {
	splitReq := strings.Split(req, sep)

	if len(splitReq) == 0 {
		return Cmd{}, errors.New("empty request")
	}

	// first element should be the command
	cmd := splitReq[0]
	op, ok := reverseOps[cmd]
	if !ok {
		return Cmd{}, fmt.Errorf("unsupported command: %s", cmd)
	}

	// the rest are all args
	args := splitReq[1:]
	if len(args) != expectedArgsCount[op] {
		return Cmd{}, fmt.Errorf(
			"expected %d arguments for command '%s', but got %d",
			expectedArgsCount[op],
			op.String(),
			len(args),
		)
	}

	return Cmd{Type: op, Args: args}, nil
}
