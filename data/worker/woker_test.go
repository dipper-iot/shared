package worker

import (
	"context"
	"gitlab.com/dipper-iot/shared/util"
	"log"
	"sync"
	"testing"
)

func TestDataWorker(t *testing.T) {
	tests := []struct {
		name       string
		dataWorker DataWorker[string]
		args       []string
		limit      int
	}{
		{
			name:       "test 1",
			dataWorker: NewDataWorker[string](10000),
			args:       []string{"1"},
			limit:      10000,
		},
		{
			name:       "test 2",
			dataWorker: NewDataWorker[string](10000),
			args:       []string{"1", "2", "3"},
			limit:      10000,
		},
		{
			name:       "test 3",
			dataWorker: NewDataWorker[string](10000),
			args:       []string{"1", "3"},
			limit:      10000,
		},
		{
			name:       "test 4",
			dataWorker: NewDataWorker[string](10000),
			args:       []string{"1", "2", "3", "4", "76756", "4545"},
			limit:      10000,
		},
		{
			name:       "test 4 limit",
			dataWorker: NewDataWorker[string](3),
			args:       []string{"1", "2", "3", "4", "76756", "4545"},
			limit:      3,
		},
		{
			name:       "test 5 limit",
			dataWorker: NewDataWorker[string](6),
			args:       []string{"1", "2", "3", "4", "76756", "4545", "df", "67"},
			limit:      6,
		},
	}
	wg := &sync.WaitGroup{}
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				defer wg.Done()
			}()
			wg2 := &sync.WaitGroup{}
			rs := make([]string, 0)
			tt.dataWorker.AddWorker(func(ctx context.Context, data string) error {
				defer wg2.Done()
				rs = append(rs, data)
				return nil
			})
			for i, arg := range tt.args {
				tt.dataWorker.Add(arg)
				if (i + 1) <= tt.limit {
					wg2.Add(1)
				}
			}
			go tt.dataWorker.Run(context.TODO())
			wg2.Wait()
			for _, item := range rs {
				if !util.StringInSlice(item, tt.args) {
					t.Errorf("Get Value Not Match")
				}
			}

			if tt.limit < len(tt.args) {

				if len(rs) != tt.limit {
					log.Println(rs)
					t.Errorf("Not Limit data")
				}

			} else {
				if len(rs) != len(tt.args) {
					t.Errorf("Get Value Not Full")
				}
			}

			tt.dataWorker.Stop()
		})
	}
	wg.Wait()
}
