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

// New constructs a key-value store.
func New() *Store {
	return &Store{
		underlying: make(map[string]string),
		m:          sync.Mutex{},
	}
}

// Get returns a value from the store.
//
// If the value exists, the second return value will return true.
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

// Stop deletes the underlying map.
func (s *Store) Stop() error {
	s.m.Lock()
	defer s.m.Unlock()

	s.underlying = nil

	return nil
}
