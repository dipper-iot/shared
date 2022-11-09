package collection

import (
	"gitlab.com/dipper-iot/shared/data/queue"
	"gitlab.com/dipper-iot/shared/distributed"
)

type DataCollection[V any] struct {
	queue *queue.QueueBase[V]
	size  int64
	lock  *distributed.TryLock
}

func NewData[V any]() *DataCollection[V] {
	return &DataCollection[V]{
		queue: queue.NewQueueBase[V](),
		lock:  distributed.NewTryLock(),
		size:  0,
	}
}

func (d *DataCollection[V]) Push(item V) bool {
	if !d.lock.Lock() {
		return false
	}
	d.queue.Push(item)
	d.size++
	d.lock.UnLock()
	return true
}

func (d *DataCollection[V]) Pop() (V, bool) {
	var (
		item V
		ok   bool
	)
	if !d.lock.Lock() {
		return item, false
	}
	item, ok = d.queue.Pop()
	d.size--
	d.lock.UnLock()
	return item, ok
}

func (d DataCollection[V]) Size() int64 {
	return d.size
}
