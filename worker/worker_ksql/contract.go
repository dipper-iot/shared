package worker_ksql

type IPayloadCollection[V any] interface {
	Init(workerNumber int, configMap map[string]interface{}) error
	Add(data V)
	Get() *Payload[V]
	HandlerPayload(data *Payload[V]) error
	HandlerError(data *Payload[V], err error)
}
