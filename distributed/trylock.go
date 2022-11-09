package distributed

type TryLock struct {
	lock bool
	c    chan struct{}
}

// tạo một TryLock
func NewTryLock() *TryLock {
	var l = &TryLock{
		c: make(chan struct{}, 1),
	}
	l.c <- struct{}{}
	return l
}

// TryLock try TryLock, trả về kết qủa TryLock là true/false
func (l TryLock) Lock() bool {
	TryLockResult := false
	select {
	case <-l.c:
		TryLockResult = true
		l.lock = true
	default:
	}
	return TryLockResult
}

// UnTryLock , giải phóng TryLock
func (l TryLock) UnLock() {
	l.c <- struct{}{}
	l.lock = false
}

func (l TryLock) IsLock() bool {
	return l.lock
}
