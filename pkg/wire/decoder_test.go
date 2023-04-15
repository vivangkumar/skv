package wire_test

import (
	"testing"

	"github.com/vivangkumar/skv/pkg/test"
	"github.com/vivangkumar/skv/pkg/wire"
)

func TestDecodeCmd(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		req := ""
		_, err := wire.DecodeCmd(req)
		if err == nil {
			t.Error("expected error not to be nil")
			t.Fail()
		}
	})

	t.Run("valid request", func(t *testing.T) {
		req := "set:key:value"

		cmd, err := wire.DecodeCmd(req)
		if err != nil {
			t.Error("expected error to be nil")
			t.Fail()
		}

		test.Equal(t, cmd.Type, wire.Set)

		test.Equal(t, 2, len(cmd.Args))
		test.Equal(t, "key", cmd.Args[0])
		test.Equal(t, "value", cmd.Args[1])
	})

	t.Run("unknown command in request", func(t *testing.T) {
		req := "abc:key:value"

		_, err := wire.DecodeCmd(req)
		if err == nil {
			t.Error("expected error, but got nil")
			t.Fail()
		}
	})

	t.Run("argument mismatch", func(t *testing.T) {
		req := "set:key:value:value"

		_, err := wire.DecodeCmd(req)
		if err == nil {
			t.Error("expected error, but got nil")
			t.Fail()
		}
	})
}
