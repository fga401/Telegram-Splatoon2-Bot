package cleaner

import (
	"container/heap"
	"time"

	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/driver/cache/internal/cleaner/internal"
)

// Cleaner cleans expired item in time.
type Cleaner struct {
	keys   *internal.ExpiredKeys
	update map[string]time.Time
	cache  cache.Cache

	taskChan chan internal.ExpiredKey
}

// New returns a new Cleaner
func New(cache cache.Cache) *Cleaner {
	cleaner := &Cleaner{
		keys:     internal.NewExpiredKeys(),
		cache:    cache,
		taskChan: make(chan internal.ExpiredKey),
	}
	heap.Init(cleaner.keys)
	go cleaner.cleanRoutine()
	return cleaner
}

// Set sets the expired key.
func (c *Cleaner) Set(key []byte, expiration time.Time) {
	c.taskChan <- internal.ExpiredKey{
		Key:        key,
		Expiration: expiration,
	}
}

func (c *Cleaner) cleanRoutine() {
	timer := time.NewTimer(0)
	// drain the channel
	// the first value was created by New
	<-timer.C
	drained := true
	for {
		select {
		case task := <-c.taskChan:
			if len(c.keys.Slice()) > 0 {
				pos := c.keys.Pos(task)
				if pos != internal.EmptyPos {
					// task already in heap
					c.keys.Slice()[pos] = task
					heap.Fix(c.keys, pos)
				}
			}
			heap.Push(c.keys, task)
			nextElem := c.keys.Slice()[0]
			// timer is running
			if !timer.Stop() && !drained {
				// timer is stopped, drain the channel
				// the channel must not be empty since read operation only happen on 'case <-c.timer.C'
				<-timer.C
			}
			timer.Reset(time.Until(nextElem.Expiration))
			drained = false
		case <-timer.C:
			drained = true
			elem := heap.Pop(c.keys).(internal.ExpiredKey)
			c.cache.Del(elem.Key)
			if len(c.keys.Slice()) > 0 {
				nextElem := c.keys.Slice()[0]
				timer.Reset(time.Until(nextElem.Expiration))
			}
		}
	}
}
