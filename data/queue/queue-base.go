package queue

// QueueBase represents a double-ended QueueBase.
// The zero value is an empty QueueBase ready to use.
/* A real stack based upon a linked list */
type QueueBase[V any] struct {
	// PushBack writes to rep[back] then increments back; PushFront
	// decrements front then writes to rep[front]; len(rep) is a power
	// of two; unused slots are nil and not garbage.
	head  *itemNode[V]
	first *itemNode[V]
}

// IntNode is a pointer data structure for holding an integer
type itemNode[V any] struct {
	value V
	left  *itemNode[V]
	right *itemNode[V]
}

func NewQueueBase[V any]() *QueueBase[V] {
	return &QueueBase[V]{}
}

// Show returns the total number of elements currently stored in the LinkedListStack
func (s QueueBase[V]) Show() []V {
	var result []V
	for current := s.head; current != nil; current = current.right {
		result = append(result, current.value)
	}
	return result
}

func (s *QueueBase[V]) Push(n V) {
	new := itemNode[V]{value: n}
	if s.head == nil {
		s.head = &new
		s.first = &new
	} else {
		new.right = s.head
		s.head.left = &new
		s.head = &new
	}
}

func (s *QueueBase[V]) Pop() (V, bool) {
	var item V

	if s.head == nil {
		return item, false
	}

	item = s.head.value
	right := s.head.right

	s.head.left = nil
	s.head.right = nil

	s.head = right
	return item, true
}

func (s *QueueBase[V]) Size() int64 {

	var size int64 = 0

	var first = s.first
	for {
		if first == nil {
			break
		}
		first = first.left
		size++
	}

	return size
}
