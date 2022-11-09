package worker

import (
	"context"
	"errors"
	"gitlab.com/dipper-iot/shared/data/collection"
	"gitlab.com/dipper-iot/shared/logger"
	"time"
)

type dataWorker[V any] struct {
	workers    chan Worker[V]
	ctx        context.Context
	cancel     context.CancelFunc
	collection *collection.Collection[V]
}

func NewDataWorker[V any](limit int64) DataWorker[V] {
	return &dataWorker[V]{
		workers:    make(chan Worker[V]),
		collection: collection.NewCollection[V](limit),
	}
}

func (d *dataWorker[V]) AddWorker(worker Worker[V]) {
	go func() {
		d.workers <- worker
	}()
}

func (d *dataWorker[V]) Add(data V) {
	d.collection.Push(data)
}

func (d *dataWorker[V]) Run(ctx context.Context) {
	d.ctx, d.cancel = context.WithCancel(ctx)
	for {
		listData := d.collection.Get(10, time.Millisecond*100)
		for _, data := range listData {
			select {
			case worker := <-d.workers:
				{
					go func() {
						defer func(d *dataWorker[V], worker Worker[V]) {
							d.workers <- worker
						}(d, worker)

						err := worker(ctx, data)
						if err == nil {
							return
						}
						logger.Error(err)
						if errors.Is(err, ErrorRollback) {
							d.Add(data)
						}
					}()
				}
			case <-d.ctx.Done():
				return
			}
		}

	}

}

func (d *dataWorker[V]) Stop() error {
	d.cancel()
	return nil
}
