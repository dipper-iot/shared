package worker_file

type Payload struct {
	Id       string                 `json:"id"`
	Version  string                 `json:"version"`
	Token    string                 `json:"token"`
	Url      string                 `json:"url"`
	Data     []byte                 `json:"data"`
	MetaData map[string]interface{} `json:"meta_data"`
}
