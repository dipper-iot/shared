package kubernetes

type Service struct {
	ServiceName      string
	Version          string
	Image            string
	Replicas         int
	Envs             []Env
	Ports            []Port
	Service          bool
	EnvFromConfigMap string
}

type Env struct {
	Name  string
	Value string
}

type Port struct {
	Name string
	Port string
}
