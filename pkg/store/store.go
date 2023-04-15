package store

import (
	"sync"
)

// Store represents a rudimentary
// string key value store.
//
// It is safe for concurrent access.
type Store struct {
	m          sync.Mutex
	underlying map[string]string
}

// NewStore constructs a key-value store.
func NewStore() *Store {
	return &Store{
		underlying: make(map[string]string),
		m:          sync.Mutex{},
	}
}

// Get returns a value with a bool for the given key.
//
// bool works according to regular go semantics.
func (s *Store) Get(k string) (string, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	v, ok := s.underlying[k]
	return v, ok
}

// Set sets the value of the key to the value.
func (s *Store) Set(k string, v string) {
	s.m.Lock()
	defer s.m.Unlock()

	s.underlying[k] = v
}

// Delete removes the key from the store.
func (s *Store) Delete(k string) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.underlying, k)
}

func (s *Store) Stop() error {
	s.m.Lock()
	defer s.m.Unlock()

	s.underlying = nil

	return nil
}
