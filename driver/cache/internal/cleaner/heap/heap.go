package heap

import (
	"telegram-splatoon2-bot/driver/cache/internal/cleaner/model"
)

// Heap is a min heap whose item's type is ExpiredKey
type Heap struct {
	keys model.ExpiredKeys
}

// New returns a new Heap given items
func New(keys model.ExpiredKeys) *Heap {
	// heapify
	n := keys.Len()
	if keys == nil {
		keys = make(model.ExpiredKeys, 0)
	}
	heap := &Heap{keys}
	for i := n/2 - 1; i >= 0; i-- {
		heap.down(i, n)
	}
	return heap
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap) Push(x *model.ExpiredKey) {
	h.keys.Push(x)
	h.up(h.keys.Len() - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *Heap) Pop() *model.ExpiredKey {
	n := h.keys.Len() - 1
	h.keys.Swap(0, n)
	h.down(0, n)
	return h.keys.Pop()
}

// Peek returns the minimum element (according to Less) from the heap, but doesn't remove it.
// If heap is empty, return nil.
// The complexity is O(1).
func (h *Heap) Peek() *model.ExpiredKey {
	if h.keys.Len() == 0 {
		return nil
	}
	return h.keys[0]
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap) Remove(i int) *model.ExpiredKey {
	n := h.keys.Len() - 1
	if n != i {
		h.keys.Swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.keys.Pop()
}

// Empty returns true if heap is empty.
// The complexity is O(1).
func (h *Heap) Empty() bool {
	return h.keys.Len() == 0
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap) Fix(i int) {
	if !h.down(i, h.keys.Len()) {
		h.up(i)
	}
}

func (h *Heap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.keys.Less(j, i) {
			break
		}
		h.keys.Swap(i, j)
		j = i
	}
}

func (h *Heap) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.keys.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.keys.Less(j, i) {
			break
		}
		h.keys.Swap(i, j)
		i = j
	}
	return i > i0
}
