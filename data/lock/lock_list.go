package lock

import "context"

type ListLocker interface {
	Waiting()
	Error() error
	Unlock()
	IsLock() bool
}

type LockList interface {
	Locker(ctx context.Context, list []string) (bool, ListLocker)
}
