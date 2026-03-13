package config

type Config struct {
	ConsulConfig ConsulConfig
	HTTP         HTTPConfig
}

type ConsulConfig struct {
	Address    string `env:"CONSUL_ADDRESS"`
	Scheme     string `env:"CONSUL_SCHEME"`
	DataCenter string `env:"CONSUL_DATA_CENTER"`
	Token      string `env:"CONSUL_TOKEN"`
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT"`
}
