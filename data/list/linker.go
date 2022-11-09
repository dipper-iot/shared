package list

import "sync"

type Linker[V any] struct {
	lock  *sync.Mutex
	start *NoteLinker[V]
	end   *NoteLinker[V]
}

func NewLinker[V any]() *Linker[V] {
	return &Linker[V]{
		end:   nil,
		start: nil,
		lock:  &sync.Mutex{},
	}
}

func (l *Linker[V]) First() *NoteLinker[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.start
}

func (l *Linker[V]) End() *NoteLinker[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.start
}

func (l *Linker[V]) Add(data V) *NoteLinker[V] {
	l.lock.Lock()
	defer l.lock.Unlock()

	newItem := &NoteLinker[V]{
		Data:  data,
		left:  l.end,
		right: nil,
	}

	if l.start == nil {
		l.start = newItem
	}

	if l.end != nil {
		l.end.right = newItem
	}

	l.end = newItem
	return newItem
}

func (l *Linker[V]) Remove(item *NoteLinker[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if item == l.start {
		l.start = item.right
		if l.start != nil {
			l.start.left = nil
		}
		return
	}

	if item == l.end {
		l.end = item.left
		if l.end != nil {
			l.end.right = nil
		}
		return
	}

	item.left = item.right
}

func (l Linker[V]) IsEnd(item *NoteLinker[V]) bool {
	if l.end == nil {
		return true
	}
	return l.end == item
}

// NoteLinker
type NoteLinker[V any] struct {
	left  *NoteLinker[V]
	right *NoteLinker[V]
	Data  V
}

func (n NoteLinker[V]) Next() (*NoteLinker[V], bool) {
	if n.right == nil {
		return nil, false
	}
	return n.right, true
}

func (n NoteLinker[V]) Prev() (*NoteLinker[V], bool) {
	if n.left == nil {
		return nil, false
	}
	return n.left, true
}
