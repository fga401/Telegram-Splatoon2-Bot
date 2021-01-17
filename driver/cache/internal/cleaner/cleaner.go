package cleaner

import (
	"sync"
	"time"

	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/driver/cache/internal/cleaner/heap"
	"telegram-splatoon2-bot/driver/cache/internal/cleaner/model"
)

// Cleaner cleans expired item in time.
type Cleaner struct {
	heap  *heap.Heap
	cache cache.Cache

	mutex    sync.Mutex
	taskChan chan *model.ExpiredKey
	timer    *time.Timer
}

// New returns a new Cleaner
func New(cache cache.Cache) *Cleaner {
	cleaner := &Cleaner{
		heap:     heap.New(nil),
		cache:    cache,
		timer:    time.NewTimer(0),
		taskChan: make(chan *model.ExpiredKey),
	}
	go cleaner.clean()
	return cleaner
}

// Set sets the expired key.
func (c *Cleaner) Set(key []byte, expiration time.Time) {
	// don't block the caller
	go func() {
		c.taskChan <- &model.ExpiredKey{
			Key:        key,
			Expiration: expiration,
		}
	}()
}

func (c *Cleaner) clean() {
	// drain the channel
	// the first value was created by New
	<-c.timer.C
	for {
		select {
		case task := <-c.taskChan:
			c.mutex.Lock()
			nextElem := c.heap.Peek()
			if nextElem == nil {
				// no task in queue, timer is already stopped
				c.timer.Reset(time.Until(task.Expiration))
			} else if task.Expiration.Unix() < nextElem.Expiration.Unix() {
				// timer is running
				if !c.timer.Stop() {
					// timer is stopped, drain the channel
					// the channel must not be empty since read operation only happen on 'case <-c.timer.C'
					<-c.timer.C
				}
				c.timer.Reset(time.Until(task.Expiration))
			}
			c.heap.Push(task)
			c.mutex.Unlock()
		case <-c.timer.C:
			c.mutex.Lock()
			elem := c.heap.Pop()
			c.cache.Del(elem.Key)
			nextElem := c.heap.Peek()
			if nextElem != nil {
				c.timer.Reset(time.Until(nextElem.Expiration))
			}
			c.mutex.Unlock()
		}
	}
}
