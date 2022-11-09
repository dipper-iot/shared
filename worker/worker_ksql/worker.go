package worker_ksql

import "time"

type Worker[V any] struct {
	dataPool   chan *Payload[V]
	collection IPayloadCollection[V]
	quit       chan bool
}

func NewWorker[V any](collection IPayloadCollection[V]) *Worker[V] {
	return &Worker[V]{
		collection: collection,
		dataPool:   make(chan *Payload[V]),
		quit:       make(chan bool),
	}
}

func (w *Worker[V]) Run(workerPool chan Worker[V]) {
	go func() {
		for {
			workerPool <- *w
			select {
			case <-w.quit:
				return
			case data := <-w.dataPool:
				{
					if data == nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					err := w.collection.HandlerPayload(data)
					if err != nil {
						w.collection.HandlerError(data, err)
					}
				}
			}
		}
	}()
}

func (w *Worker[V]) AddJob(data *Payload[V]) {
	w.dataPool <- data
}

func (w *Worker[V]) Stop() {
	w.quit <- true
}
