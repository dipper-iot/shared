package worker_ksql

type Payload[V any] struct {
	Data []V `json:"data"`
}
