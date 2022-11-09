package queue

import (
	"sync"
)

// QueueLock represents a double-ended QueueLock.
// The zero value is an empty QueueLock ready to use.
/* A real stack based upon a linked list */
type QueueLock[V any] struct {
	*QueueBase[V]
	lock *sync.Mutex
}

func NewQueue[V any]() Queue[V] {
	return &QueueLock[V]{
		QueueBase: NewQueueBase[V](),
		lock:      new(sync.Mutex),
	}
}

func (s *QueueLock[V]) Push(data V) {
	s.lock.Lock()
	s.QueueBase.Push(data)
	s.lock.Unlock()
}

func (s *QueueLock[V]) Pop() (V, bool) {
	var (
		item V
		ok   bool
	)
	s.lock.Lock()
	item, ok = s.QueueBase.Pop()
	s.lock.Unlock()
	return item, ok
}
