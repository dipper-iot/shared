package worker_file

import (
	"context"
	"encoding/json"
	"gitlab.com/dipper-iot/shared/logger"
	"gitlab.com/dipper-iot/shared/util"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type FilePayloadCollection struct {
	DataPool    chan Payload
	WorkerPool  chan Worker
	FolderSetup string
	Folder      string
	ctx         context.Context
	sync.Mutex
}

func NewFilePayloadCollection(ctx context.Context, maxDataPool int) *FilePayloadCollection {
	return &FilePayloadCollection{
		ctx:        ctx,
		DataPool:   make(chan Payload, maxDataPool),
		WorkerPool: make(chan Worker),
		Mutex:      sync.Mutex{},
	}
}

//2006-01-02T15:04:05Z07:00
func (f *FilePayloadCollection) setFolderWithTime() {
	f.Folder = filepath.Join(f.FolderSetup, time.Now().Format("2006_01_02_15_04_05"))
}

func (f *FilePayloadCollection) Init(workerNumber int, configMap map[string]interface{}) error {

	f.FolderSetup = configMap["folder"].(string)
	f.setFolderWithTime()
	go func() {
		tinker := time.Tick(5 * time.Minute)
		for {
			select {
			case <-tinker:
				{
					f.setFolderWithTime()
				}
			case <-f.ctx.Done():
				{
					return
				}
			}
		}
	}()

	for i := 0; i < workerNumber; i++ {
		worker := NewWorker(f)
		worker.Run(f.WorkerPool)
	}

	go f.dispatch()
	return nil
}

func (f *FilePayloadCollection) dispatch() {
	for {
		select {
		case data := <-f.Get():
			{
				// a job request has been received
				go func(data Payload) {
					// try to obtain a worker_broker job channel that is available.
					// this will block until a worker_broker is idle
					worker := <-f.WorkerPool

					// dispatch the job to the worker_broker job channel
					worker.AddJob(data)
				}(data)
			}
		case <-f.ctx.Done():
			{
				return
			}

		}
	}
}

func (f *FilePayloadCollection) Add(data Payload) {
	f.DataPool <- data
}

func (f *FilePayloadCollection) Get() <-chan Payload {
	return f.DataPool
}

func (f *FilePayloadCollection) HandlerPayload(data Payload) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return NewErrorPayload("json", err)
	}
	util.EnsureDir(f.Folder)
	err = util.SaveFile(path.Join(f.Folder, data.Id), bytes)
	if err != nil {
		return NewErrorPayload("save", err)
	}
	return nil
}

func (f *FilePayloadCollection) HandlerError(data Payload, err error) {
	errPay, _ := ConvertError(err)
	if errPay == nil {
		logger.Error(err)
		return
	}
	if errPay.Code == "command" || errPay.Code == "save" {
		go func() {
			f.DataPool <- data
		}()
	}
	logger.Error(err)
}
