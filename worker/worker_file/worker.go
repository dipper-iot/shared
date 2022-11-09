package worker_file

type Worker struct {
	dataPool   chan Payload
	collection IPayloadCollection
	quit       chan bool
}

func NewWorker(collection IPayloadCollection) *Worker {
	return &Worker{
		collection: collection,
		dataPool:   make(chan Payload),
		quit:       make(chan bool),
	}
}

func (w *Worker) Run(workerPool chan Worker) {
	go func() {
		for {
			workerPool <- *w
			select {
			case <-w.quit:
				return
			case data := <-w.dataPool:
				{
					err := w.collection.HandlerPayload(data)
					if err != nil {
						w.collection.HandlerError(data, err)
					}
				}
			}
		}
	}()
}

func (w *Worker) AddJob(data Payload) {
	w.dataPool <- data
}

func (w *Worker) Stop() {
	w.quit <- true
}
