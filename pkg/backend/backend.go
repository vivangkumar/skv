package backend

import (
	"context"
)

type store interface {
	Set(k string, v string)
	Get(k string) (string, bool)
	Delete(k string)
	Stop() error
}

// Backend represents an abstraction to Store
// It exists to add a layer of context cancellation.
type Backend struct {
	st store
}

// NewBackend creates a bakcend for the given store
func NewBackend(st store) *Backend {
	return &Backend{st: st}
}

// Set sets a value to the store while allowing
// context cancellation.
func (b *Backend) Set(ctx context.Context, cmd SetCmd) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		b.st.Set(cmd.Key, cmd.Value)
	}

	return nil
}

// Get returns a value from the store.
//
// If no value is found, an empty string is returned.
func (b *Backend) Get(ctx context.Context, cmd GetCmd) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		v, ok := b.st.Get(cmd.Key)
		if !ok {
			return "", nil
		}

		return v, nil
	}
}

// Delete removes a key already set in the store.
func (b *Backend) Delete(ctx context.Context, cmd DelCmd) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		b.st.Delete(cmd.Key)
	}

	return nil
}

func (b *Backend) Stop() error {
	return b.st.Stop()
}
