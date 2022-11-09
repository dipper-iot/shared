package queue

// Queue represents a double-ended queue.
// The zero value is an empty queue ready to use.
/* A real stack based upon a linked list */
type Queue[V any] interface {
	Push(data V)
	Size() int64
	Pop() (V, bool)
}
