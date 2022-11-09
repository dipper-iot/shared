package worker

import "errors"

var (
	ErrorNotFound = errors.New("WORKER_DATA_NOT_FOUND")
	ErrorRollback = errors.New("WORKER_ROLLBACK")
	ErrorNotRun   = errors.New("WORKER_NOT_RUN")
)
