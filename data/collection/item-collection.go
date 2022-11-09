package collection

type itemCollection[V any] struct {
	left  *itemCollection[V]
	right *itemCollection[V]
	value *DataCollection[V]
}

func newItemCollection[V any](left *itemCollection[V], right *itemCollection[V], value *DataCollection[V]) *itemCollection[V] {

	return &itemCollection[V]{
		left:  left,
		right: right,
		value: value,
	}
}

func newItemCollectionOnly[V any](value *DataCollection[V]) *itemCollection[V] {

	return &itemCollection[V]{
		left:  nil,
		right: nil,
		value: value,
	}
}

func newItemCollectionBlank[V any](left *itemCollection[V], right *itemCollection[V]) *itemCollection[V] {
	data := NewData[V]()

	return &itemCollection[V]{
		left:  left,
		right: right,
		value: data,
	}
}

func newItemCollectionLeftBlank[V any](left *itemCollection[V]) *itemCollection[V] {
	data := NewData[V]()

	return &itemCollection[V]{
		left:  left,
		right: nil,
		value: data,
	}
}
