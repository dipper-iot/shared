package worker_file

type IPayloadCollection interface {
	Init(workerNumber int, configMap map[string]interface{}) error
	Add(data Payload)
	Get() <-chan Payload
	HandlerPayload(data Payload) error
	HandlerError(data Payload, err error)
}
