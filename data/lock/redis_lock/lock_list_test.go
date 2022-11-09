package redis_lock

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.com/dipper-iot/shared/load/rs"
	"gitlab.com/dipper-iot/shared/logger"
	"golang.org/x/net/context"
	"sync"
	"testing"
	"time"
)

var (
	client *redis.Client
)

func init() {
	var err error
	client, err = rs.NewRedisToEnv()
	if err != nil {
		logger.Error(err)
	}

	//logger.Init(logger.WithLevel(logger.TraceLevel))
}

func nilKey(lockKey string) bool {
	re, err := client.Get(context.TODO(), lockKey).Result()
	if redis.Nil == err {
		return true
	}
	logger.Info(re)
	return false
}

func TestLockListUnLock(t *testing.T) {
	lockList := NewLockListRedis(client, 1*time.Nanosecond)
	lockList.Run()

	listData := [][]string{
		{"l1", "l2"},
		{"l3"},
		{"l4"},
		{"l5"},
		{"l6"},
	}
	result := map[string]int{
		"l1": 1,
		"l2": 1,
		"l3": 1,
		"l4": 1,
		"l5": 1,
		"l6": 1,
	}
	counts := sync.Map{}

	wg := &sync.WaitGroup{}

	for index, list := range listData {
		wg.Add(1)
		go func(index int, list []string) {
			defer wg.Done()
			_, locker := lockList.Locker(context.Background(), list)
			locker.Waiting()
			for _, item := range list {
				countData, success := counts.Load(item)
				var count int = 1
				if success {
					count = countData.(int)
					counts.Store(item, count+1)
				} else {
					counts.Store(item, count)
				}

			}
			logger.Infof("Success list %d %s", index, list)
			viewMap(counts)
			locker.Unlock()
		}(index, list)

	}

	wg.Wait()

	for key, count := range result {
		countData, success := counts.Load(key)
		if !success {
			t.Errorf("Not match count with %s result is %d but actually nil", key, count)
			continue
		}
		resultCount := countData.(int)
		if resultCount != count {
			t.Errorf("Not match count with %s result is %d but actually %d", key, count, resultCount)
		}
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

}

func TestLockList(t *testing.T) {
	lockList := NewLockListRedis(client, 1*time.Nanosecond)
	lockList.Run()

	listData := [][]string{
		{"l1", "l2", "l3"},
		{"l1", "l3"},
		{"l3", "l2"},
		{"l4"},
		{"l2"},
		{"l3", "l4"},
	}
	result := map[string]int{
		"l1": 2,
		"l2": 3,
		"l3": 4,
		"l4": 2,
	}
	counts := sync.Map{}

	wg := &sync.WaitGroup{}

	for index, list := range listData {
		wg.Add(1)
		go func(index int, list []string) {
			defer wg.Done()
			_, locker := lockList.Locker(context.Background(), list)
			locker.Waiting()
			for _, item := range list {
				countData, success := counts.Load(item)
				var count int = 1
				if success {
					count = countData.(int)
					counts.Store(item, count+1)
				} else {
					counts.Store(item, count)
				}

			}
			logger.Infof("Success list %d %s", index, list)
			viewMap(counts)
			locker.Unlock()
		}(index, list)

	}

	wg.Wait()

	for key, count := range result {
		countData, success := counts.Load(key)
		if !success {
			t.Errorf("Not match count with %s result is %d but actually nil", key, count)
			continue
		}
		resultCount := countData.(int)
		if resultCount != count {
			t.Errorf("Not match count with %s result is %d but actually %d", key, count, resultCount)
		}
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

}

func viewMap(mapData sync.Map) {
	str := ""
	mapData.Range(func(key, value interface{}) bool {
		count := value.(int)
		str = fmt.Sprintf("%s %s=%d", str, key, count)
		return true
	})
	logger.Infof("map[%s ]", str)
}

func TestLockList2(t *testing.T) {
	lockList := NewLockListRedis(client, 1*time.Nanosecond)
	lockList.Run()

	listData := [][]string{
		{"l1", "l2", "l10"},
		{"l1", "l3"},
		{"l3", "l2"},
		{"l4"},
		{"l2"},
		{"l3", "l4"},
		{"l5"},
		{"l3", "l1", "l5", "l7"},
		{"l7"},
		{"l8"},
		{"l9"},
		{"l9", "l10"},
	}
	result := map[string]int{
		"l1":  3,
		"l2":  3,
		"l3":  4,
		"l4":  2,
		"l5":  2,
		"l7":  2,
		"l8":  1,
		"l9":  2,
		"l10": 2,
	}
	counts := sync.Map{}

	wg := &sync.WaitGroup{}

	for index, list := range listData {
		wg.Add(1)
		go func(index int, list []string) {
			defer wg.Done()
			_, locker := lockList.Locker(context.Background(), list)
			locker.Waiting()
			defer locker.Unlock()
			for _, item := range list {
				countData, success := counts.Load(item)
				var count int = 1
				if success {
					count = countData.(int)
					counts.Store(item, count+1)
				} else {
					counts.Store(item, count)
				}

			}
			logger.Infof("Success list %d %s", index, list)
			viewMap(counts)
		}(index, list)

	}

	wg.Wait()

	for key, count := range result {
		countData, success := counts.Load(key)
		if !success {
			t.Errorf("Not match count with %s result is %d but actually nil", key, count)
			continue
		}
		resultCount := countData.(int)
		if resultCount != count {
			t.Errorf("Not match count with %s result is %d but actually %d", key, count, resultCount)
		}
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

}

func inc(list []string, counts *sync.Map) {
	for _, item := range list {
		countData, success := counts.Load(item)
		var count int = 1
		if success {
			count = countData.(int)
			counts.Store(item, count+1)
		} else {
			counts.Store(item, count)
		}

	}
}

func TestLockListCancel(t *testing.T) {
	lockList := NewLockListRedis(client, 1*time.Nanosecond)
	lockList.Run()

	listData := [][]string{
		{"l1", "l2", "l3"},
		{"l1", "l3"},
		{"l3", "l2"},
		{"l4"},
		{"l2"},
		{"l3", "l4"},
		{"l1"},
	}
	result := map[string]int{
		"l1": 1,
		"l2": 2,
		"l3": 2,
		"l4": 2,
	}
	counts := sync.Map{}

	ctx1, _ := context.WithTimeout(context.TODO(), 1*time.Nanosecond)
	ctx2, _ := context.WithTimeout(context.TODO(), 20*time.Second)

	wg := &sync.WaitGroup{}

	for index, list := range listData {
		wg.Add(1)

		if index < 2 {
			go func(index int, list []string) {
				defer wg.Done()
				_, locker := lockList.Locker(ctx1, list)
				locker.Waiting()
				defer locker.Unlock()
				if locker.Error() == context.Canceled || locker.Error() == context.DeadlineExceeded {
					return
				}
				inc(list, &counts)
				logger.Infof("Success list %d %s", index, list)
				viewMap(counts)
				locker.Unlock()
			}(index, list)
		} else {
			go func(index int, list []string) {
				defer wg.Done()
				_, locker := lockList.Locker(ctx2, list)
				locker.Waiting()
				defer locker.Unlock()
				if locker.Error() == context.Canceled || locker.Error() == context.DeadlineExceeded {
					return
				}
				inc(list, &counts)
				logger.Infof("Success list %d %s", index, list)
				viewMap(counts)
				locker.Unlock()
			}(index, list)
		}

	}

	wg.Wait()

	for key, count := range result {
		countData, success := counts.Load(key)
		if !success {
			t.Errorf("Not match count with %s result is %d but actually nil", key, count)
			continue
		}
		resultCount := countData.(int)
		if resultCount != count {
			t.Errorf("Not match count with %s result is %d but actually %d", key, count, resultCount)
		}
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

}

func TestLockListCancelUnlock(t *testing.T) {
	lockList := NewLockListRedis(client, 1*time.Nanosecond)
	lockList.Run()

	list := []string{"l1", "l2", "l3"}

	ctx1, _ := context.WithTimeout(context.TODO(), 1*time.Nanosecond)

	_, locker := lockList.Locker(ctx1, list)
	locker.Waiting()
	defer locker.Unlock()

	for _, key := range list {
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

}

func TestLockListCancelUnlockAndLock(t *testing.T) {
	lockList := NewLockListRedis(client, 1*time.Nanosecond)
	lockList.Run()

	list := []string{"l1", "l2", "l3"}
	list2 := []string{"l1", "l2", "l4"}

	ctx1, _ := context.WithTimeout(context.TODO(), 1*time.Nanosecond)
	ctx2, _ := context.WithTimeout(context.TODO(), 20*time.Second)

	go func() {
		_, locker := lockList.Locker(ctx1, list)
		locker.Waiting()
		if nil == locker.Error() {
			t.Errorf("Error is not nil")
		}
		locker.Unlock()
	}()

	go func() {
		_, locker := lockList.Locker(ctx2, list2)
		locker.Waiting()
		if nil != locker.Error() {
			t.Errorf("Error not nil")
		}
		locker.Unlock()
	}()

	for _, key := range list {
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

	for _, key := range list2 {
		if !nilKey(key) {
			t.Errorf("Key Redis not nil %s", key)
		}
	}

}
