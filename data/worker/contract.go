package worker

import "context"

type Worker[V any] func(ctx context.Context, data V) error

type DataWorker[V any] interface {
	AddWorker(worker Worker[V])
	Add(data V)
	Run(ctx context.Context)
	Stop() error
}
