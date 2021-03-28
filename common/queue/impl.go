package queue

import (
	"container/list"
	"sync"
)

// New returns a unlimited buffer Queue implemented by list.
func New() Queue {
	ret := &impl{
		cond:     sync.NewCond(&sync.Mutex{}),
		list:     list.New(),
		closer:   make(chan struct{}),
		inChan:   make(chan interface{}),
		outChan:  make(chan interface{}),
		initOnce: &sync.Once{},
	}
	go ret.routine()
	return ret
}

type impl struct {
	cond     *sync.Cond
	list     *list.List
	closer   chan struct{}
	inChan   chan interface{}
	outChan  chan interface{}
	initOnce *sync.Once
}

func (q *impl) EnqueueChan() chan<- interface{} {
	return q.inChan
}

func (q *impl) DequeueChan() <-chan interface{} {
	return q.outChan
}

func (q *impl) routine() {
	curVal := func() interface{} {
		if q.list.Len() > 0 {
			return q.list.Front().Value
		}
		return nil
	}
	outChan := func() chan interface{} {
		if q.list.Len() > 0 {
			return q.outChan
		}
		return nil
	}
	for q.list.Len() > 0 || q.inChan != nil {
		select {
		case v, ok := <-q.inChan:
			if ok {
				q.list.PushBack(v)
			} else {
				q.inChan = nil
			}
		case outChan() <- curVal():
			q.list.Remove(q.list.Front())
		}
	}
	close(q.outChan)
}
