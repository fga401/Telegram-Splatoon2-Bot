package queue

// Queue is an unlimited buffer channel
type Queue interface {
	// EnqueueChan return a write-only channel to enqueue elements.
	// Once the channel closed, the queue will be closed after all elements dequeued.
	// Close a closed enqueue channel will cause panic.
	EnqueueChan() chan<- interface{}
	// EnqueueChan return a read-only channel to dequeue elements.
	// This channel will be closed after queue closed.
	DequeueChan() <-chan interface{}
}

