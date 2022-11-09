package postgres

type Config struct {
	Addr     string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
}