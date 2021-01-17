package model

import "time"

// ExpiredKey stores key and its expiration time
type ExpiredKey struct {
	Key        []byte
	Expiration time.Time // The time when key will expire in
}

// ExpiredKeys is slice of ExpiredKey
type ExpiredKeys []*ExpiredKey

// Len is part of sort.Interface.
func (q *ExpiredKeys) Len() int {
	return len(*q)
}

// Swap is part of sort.Interface.
func (q *ExpiredKeys) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

// Less is part of sort.Interface.
func (q *ExpiredKeys) Less(i, j int) bool {
	return (*q)[i].Expiration.Unix() < (*q)[j].Expiration.Unix()
}

// Push adds x as element Len()
func (q *ExpiredKeys) Push(x *ExpiredKey) {
	*q = append(*q, x)
}

// Push removes and return element Len() - 1
func (q *ExpiredKeys) Pop() *ExpiredKey {
	n := q.Len() - 1
	ret := (*q)[n]
	*q = (*q)[0:n]
	return ret
}
