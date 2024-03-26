package wire_test

import (
	"testing"

	"github.com/vivangkumar/skv/internal/test"
	"github.com/vivangkumar/skv/internal/wire"
)

func TestEncodeReply(t *testing.T) {
	t.Run("encode ok (without args)", func(t *testing.T) {
		encoded, err := wire.EncodeReply(wire.CmdReply{
			Type: wire.OK,
			Args: []string{},
		})
		if err != nil {
			t.Error("expected error to be nil")
			t.Fail()
		}

		test.Equal(t, string(encoded), "ok\r\n")
	})

	t.Run("encode ok (with args)", func(t *testing.T) {
		encoded, err := wire.EncodeReply(wire.CmdReply{
			Type: wire.OK,
			Args: []string{"a"},
		})
		if err != nil {
			t.Error("expected error to be nil")
			t.Fail()
		}

		test.Equal(t, string(encoded), "ok:a\r\n")
	})

	t.Run("encode err", func(t *testing.T) {
		encoded, err := wire.EncodeReply(wire.CmdReply{
			Type: wire.Err,
			Args: []string{"cmd_error", "command failed"},
		})
		if err != nil {
			t.Error("expected error to be nil")
			t.Fail()
		}

		test.Equal(t, string(encoded), "err:cmd_error:command failed\r\n")
	})
}
