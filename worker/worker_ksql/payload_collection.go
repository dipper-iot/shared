package worker_ksql

import (
	"context"
	"github.com/rmoff/ksqldb-go"
	"gitlab.com/dipper-iot/shared/data/collection"
	"gitlab.com/dipper-iot/shared/ksql"
	"gitlab.com/dipper-iot/shared/logger"
	"strings"
	"time"
)

type KsqlPayloadCollection[V any] struct {
	Pool       *collection.Collection[V]
	WorkerPool chan Worker[V]
	client     *ksqldb.Client
	ctx        context.Context
}

func NewKsqlPayloadCollection[V any](ctx context.Context, client *ksqldb.Client) *KsqlPayloadCollection[V] {
	return &KsqlPayloadCollection[V]{
		ctx:        ctx,
		Pool:       collection.NewCollection[V](50),
		client:     client,
		WorkerPool: make(chan Worker[V]),
	}
}

func (f *KsqlPayloadCollection[V]) Init(workerNumber int, configMap map[string]interface{}) error {

	for i := 0; i < workerNumber; i++ {
		worker := NewWorker[V](f)
		worker.Run(f.WorkerPool)
	}

	go f.dispatch()
	return nil
}

func (f *KsqlPayloadCollection[V]) dispatch() {
	for {
		select {
		case worker := <-f.WorkerPool:
			{
				data := f.Get()
				go worker.AddJob(data)
			}
		case <-f.ctx.Done():
			{
				return
			}
		}
	}
}

func (f *KsqlPayloadCollection[V]) Add(data V) {
	f.Pool.Push(data)
}

func (f *KsqlPayloadCollection[V]) Get() *Payload[V] {
	data := f.Pool.Get(100, 500*time.Millisecond)
	if len(data) == 0 {
		return nil
	}
	return &Payload[V]{
		Data: data,
	}
}

func (f *KsqlPayloadCollection[V]) HandlerPayload(data *Payload[V]) error {

	query := make([]string, 0)

	for _, item := range data.Data {
		sql, err := ksql.QueryInsert(item)
		if err != nil {
			logger.Info(item)
			logger.Error(err)
			continue
		}
		query = append(query, sql)
	}
	if len(query) == 0 {
		return nil
	}
	err := f.client.Execute(strings.Join(query, ""))
	if err != nil {
		return NewErrorPayload("command", err)
	}
	return nil
}

func (f *KsqlPayloadCollection[V]) HandlerError(data *Payload[V], err error) {
	errPay, _ := ConvertError(err)
	if errPay == nil {
		logger.Error(err)
		return
	}
	if errPay.Code == "command" || errPay.Code == "save" {
		go func() {
			for _, datum := range data.Data {
				f.Add(datum)
			}
		}()
	}
	logger.Error(err)
}
