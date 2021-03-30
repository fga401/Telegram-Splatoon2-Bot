package queue

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

const (
	n = 1000
)

func TestQueue(t *testing.T) {
	q := New()
	in := q.EnqueueChan()
	out := q.DequeueChan()

	go func() {
		for i := 0; i < n; i++ {
			in <- i
		}
	}()
	slice := make([]int, 0, n)
	for i := 0; i < n; i++ {
		v := (<-out).(int)
		slice = append(slice, v)
	}
	for i := 0; i < n; i++ {
		require.Equal(t, i, slice[i])
	}
}
