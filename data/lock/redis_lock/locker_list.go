package redis_lock

import (
	"context"
	"gitlab.com/dipper-iot/shared/data/list"
	"gitlab.com/dipper-iot/shared/logger"
)

type LockerListRedis struct {
	ctx       context.Context
	cancel    context.CancelFunc
	lockeList *LockListRedis
	lock      chan bool
	list      []string
	isLock    bool
	isUnlock  bool
	node      *list.NoteLinker[*LockerListRedis]
	err       error
}

func (l *LockerListRedis) Error() error {
	return l.err
}

func NewLockerListRedis(ctx context.Context, lockeList *LockListRedis, list []string) *LockerListRedis {
	ctxNew, cancel := context.WithCancel(ctx)
	locker := &LockerListRedis{
		lock:      make(chan bool),
		isLock:    true,
		isUnlock:  false,
		list:      list,
		lockeList: lockeList,
		ctx:       ctxNew,
		cancel:    cancel,
		err:       nil,
	}

	return locker
}

func (l *LockerListRedis) checkCancel() {
	select {
	case <-l.ctx.Done():
		{
			l.Unlock()
			if !l.isUnlock {
				l.lockeList.removeLocker(l)
			}
		}
	}
}

func (l LockerListRedis) Waiting() {
	<-l.lock
	close(l.lock)
	logger.Tracef("Exits Waiting %s", l.list)
}

func (l *LockerListRedis) Unlock() {
	if l.isUnlock || l.isLock {
		return
	}
	err := l.lockeList.unLockLocker(context.Background(), l)
	if err != nil {
		l.err = err
	}
	l.isUnlock = true
	l.cancel()
}

func (l LockerListRedis) IsLock() bool {
	return l.isLock
}
