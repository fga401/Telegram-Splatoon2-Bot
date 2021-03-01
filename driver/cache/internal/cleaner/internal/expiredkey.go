package internal

import (
	"time"
)

const (
	// EmptyPos is the value ExpiredKeys.Pos returned when key not found.
	EmptyPos = -1
)

// ExpiredKey stores key and its expiration time
type ExpiredKey struct {
	Key        []byte
	Expiration time.Time // The time when key will expire in
}

// ExpiredKeys is slice of ExpiredKey
type ExpiredKeys struct {
	keys []ExpiredKey
	pos  map[string]int
}

// NewExpiredKeys returns ExpiredKeys object.
func NewExpiredKeys() *ExpiredKeys {
	return &ExpiredKeys{
		keys: make([]ExpiredKey, 0),
		pos:  make(map[string]int),
	}
}

// Len is part of sort.Interface.
func (q *ExpiredKeys) Len() int {
	return len(q.keys)
}

// Swap is part of sort.Interface.
func (q *ExpiredKeys) Swap(i, j int) {
	q.pos[q.keys[i].toKey()], q.pos[q.keys[j].toKey()] = j, i
	q.keys[i], q.keys[j] = q.keys[j], q.keys[i]
}

// Less is part of sort.Interface.
func (q *ExpiredKeys) Less(i, j int) bool {
	return q.keys[i].Expiration.Unix() < q.keys[j].Expiration.Unix()
}

// Push adds x as element Len()
func (q *ExpiredKeys) Push(x interface{}) {
	xx := x.(ExpiredKey)
	q.pos[xx.toKey()] = q.Len()
	q.keys = append(q.keys, xx)
}

// Pop removes and returns element Len() - 1
func (q *ExpiredKeys) Pop() interface{} {
	n := q.Len() - 1
	ret := q.keys[n]
	delete(q.pos, ret.toKey())
	q.keys = q.keys[0:n]
	return ret
}

// Pos returns the position of the ExpiredKey. If not existed, return -1.
func (q *ExpiredKeys) Pos(x ExpiredKey) int {
	if v, ok := q.pos[x.toKey()]; ok {
		return v
	}
	return EmptyPos
}

// Slice returns the ExpiredKey slice. It may be inconsistent after calling Push or Pop.
func (q *ExpiredKeys) Slice() []ExpiredKey {
	return q.keys
}

func (q *ExpiredKey) toKey() string {
	return string(q.Key)
}
