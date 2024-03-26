package store_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/vivangkumar/skv/internal/store"
	"github.com/vivangkumar/skv/internal/test"
)

func TestStore(t *testing.T) {
	s := store.New()

	t.Run("set/get", func(t *testing.T) {
		k := "key"
		expected := "value"

		s.Set(k, expected)

		actual, ok := s.Get(k)
		test.Equal(t, true, ok)
		test.Equal(t, expected, actual)
	})

	t.Run("delete", func(t *testing.T) {
		k := "mykey"

		s.Set(k, "value")
		s.Delete(k)

		_, ok := s.Get(k)
		test.Equal(t, false, ok)
	})

	t.Run("concurrent read/ write", func(t *testing.T) {
		s := store.New()

		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("k:%d", i)
				value := fmt.Sprintf("v:%d", i)

				s.Set(key, value)
				s.Get(key)
			}(i)
		}

		wg.Wait()
	})
}
