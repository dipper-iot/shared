package collection

import "C"
import (
	"context"
	"gitlab.com/dipper-iot/shared/distributed"
	"sync"
	"time"
)

type Collection[V any] struct {
	head    *itemCollection[V]
	end     *itemCollection[V]
	lock    *sync.Mutex
	tryLock *distributed.TryLock
	limit   int64
	size    int64
}

func NewCollection[V any](limit int64) *Collection[V] {
	return &Collection[V]{
		limit:   limit,
		tryLock: distributed.NewTryLock(),
		lock:    new(sync.Mutex),
	}
}

func (c *Collection[V]) Push(value V) {

	if c.head == nil {
		c.lock.Lock()
		if c.head == nil {
			data := NewData[V]()
			data.Push(value)
			item := newItemCollectionOnly[V](data)
			c.head = item
			c.end = item
			c.lock.Unlock()
			return
		}
		c.lock.Unlock()
	}

PushData:
	current := c.head
	last := c.head
	if c.Size() >= c.limit {
		return
	}
	for current != nil {
		if current.value.Push(value) {
			return
		}
		last = current
		current = current.right
	}
	if last == c.end && c.tryLock.Lock() {
		item := newItemCollectionLeftBlank[V](c.end)
		c.end = item
		c.tryLock.UnLock()
	}

	goto PushData
}

func (c *Collection[V]) Pop() (V, bool) {

	var (
		item V
		ok   bool
	)

	current := c.head
	for current != nil {
		item, ok = current.value.Pop()
		if ok {
			return item, true
		}
		current = c.head.right
	}

	return item, false
}

func (c *Collection[V]) Size() int64 {

	var size int64 = 0

	var first = c.head
	for {
		if first == nil {
			break
		}
		size = size + first.value.Size()
		first = first.right
	}

	return size
}

func (c *Collection[V]) Get(size int64, timeout time.Duration) []V {

	cxt, _ := context.WithTimeout(context.TODO(), timeout)

	list := make([]V, 0)
	var lent int64 = 0

	for {
		select {
		case <-cxt.Done():
			{
				return list
			}
		default:
			item, ok := c.Pop()
			if ok {
				lent++
				list = append(list, item)
			} else {
				time.Sleep(10)
			}
			if size == lent {
				return list
			}
		}
	}

	return list
}
