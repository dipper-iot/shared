package group

import (
	"context"
	"sync"
	"time"
)

type GroupData[V any] struct {
	queue  chan V
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewGroupData[V any](ctx context.Context) *GroupData[V] {
	ctx, cancel := context.WithCancel(ctx)
	return &GroupData[V]{
		mu:     new(sync.Mutex),
		queue:  make(chan V),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *GroupData[V]) Add(v V) {
	go func() {
		select {
		case g.queue <- v:
			return
		case <-g.ctx.Done():
			return
		}
	}()
}

func (g *GroupData[V]) Stop() {
	g.cancel()
	close(g.queue)
}

func (g *GroupData[V]) Get(timeGet time.Duration) []V {
	data := make([]V, 0)
	ctx, _ := context.WithTimeout(context.TODO(), timeGet)

LOOP:
	for {
		select {
		case <-ctx.Done():
			{
				break LOOP
			}
		case value := <-g.queue:
			{
				data = append(data, value)
				continue
			}

		}
	}

	return data
}
