package redis_lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab.com/dipper-iot/shared/data/list"
	"gitlab.com/dipper-iot/shared/data/lock"
	"gitlab.com/dipper-iot/shared/logger"
	"time"
)

type LockListRedis struct {
	timeLoop time.Duration
	linker   *list.Linker[*LockerListRedis]
	client   *redis.Client
	run      bool
}

func NewLockListRedis(client *redis.Client, timeLoop time.Duration) *LockListRedis {
	obj := &LockListRedis{
		client:   client,
		timeLoop: timeLoop,
		linker:   list.NewLinker[*LockerListRedis](),
		run:      true,
	}

	return obj
}

func (l *LockListRedis) Run() {
	go l.scan()
}
func (l *LockListRedis) Stop() {
	l.run = false
}

func (l *LockListRedis) Locker(ctx context.Context, list []string) (bool, lock.ListLocker) {
	locker := NewLockerListRedis(ctx, l, list)
	success, err := l.lockList(ctx, list)
	if err != nil {
		logger.Error(err)
	}
	if success || err != nil {
		locker.err = err
		if err != nil {
			locker.isLock = false
		}
		go func() {
			locker.lock <- true
		}()
	}
	node := l.linker.Add(locker)
	locker.node = node

	go locker.checkCancel()

	return success, locker
}

type ListItem struct {
	run     bool
	blocker *LockerListRedis
	list    []string
}

func (l *LockListRedis) scan() {
	var success bool
	for l.run {
		time.Sleep(l.timeLoop)
		item := l.linker.First()
	LockCheck:
		for {
			if item == nil {
				break LockCheck
			}

			locker := item.Data

			if locker.IsLock() {
				success, err := l.lockList(locker.ctx, locker.list)
				if success || err != nil {
					locker.err = err
					if err != nil {
						locker.isLock = false
					}
					go func() {
						locker.lock <- true
					}()
				}
			}

			if l.linker.IsEnd(item) {
				break LockCheck
			}
			item, success = item.Next()
			if !success {
				break LockCheck
			}
		}
	}
}

func (l *LockListRedis) unLockLocker(ctx context.Context, locker *LockerListRedis) error {

	err := l.unlockList(ctx, locker.list)
	if err != nil {
		logger.Error(err)
	}

	l.removeLocker(locker)

	return err
}

func (l *LockListRedis) removeLocker(locker *LockerListRedis) {
	l.linker.Remove(locker.node)
}

func (l *LockListRedis) lockList(ctx context.Context, list []string) (bool, error) {
	listSuccess := make([]string, 0)
	var (
		err         error
		lockSuccess bool
	)
	result := true
	for _, name := range list {
		resp := l.client.SetNX(ctx, name, 1, time.Second*5)
		lockSuccess, err = resp.Result()
		if err != nil {
			logger.Errorf("lock failed %s %s %s ", name, list, err)
			result = false
			break
		}
		if !lockSuccess {
			result = false
			break
		}
		//logger.Trace("lock success ", name)
		listSuccess = append(listSuccess, name)
	}

	if !result {
		if len(listSuccess) > 0 {
			err := l.unlockList(context.TODO(), listSuccess)
			if err != nil {
				logger.Errorf("lock failed %s ", err)
			}
		}

	} else {
		logger.Trace("lock success list ", list)
	}

	return result, err
}

func (l *LockListRedis) unlockList(ctx context.Context, list []string) error {
	delResp := l.client.Del(ctx, list...)
	_, err := delResp.Result()
	if err != nil {
		logger.Error("unlock failed ", err)
		return err
	}
	logger.Trace("Unlock success list ", list)
	return nil
}
