package config

type Config struct {
	RedisConfig RedisConfig
	HttpConfig  HttpConfig
	DBConfig    DBConfig
}

type RedisConfig struct {
	REDIS_CONNECTION_URL         string `env:"REDIS_CONNECTION_URL"`
	REDIS_POOL_SIZE              int    `env:"REDIS_POOL_SIZE"`
	REDIS_DIAL_TIMEOUT_SECONDS   int    `env:"REDIS_DIAL_TIMEOUT_SECONDS"`
	REDIS_READ_TIMEOUT_SECONDS   int    `env:"REDIS_READ_TIMEOUT_SECONDS"`
	REDIS_WRITE_TIMEOUT_SECONDS  int    `env:"REDIS_WRITE_TIMEOUT_SECONDS"`
	REDIS_IDLE_TIMEOUT_SECONDS   int    `env:"REDIS_IDLE_TIMEOUT_SECONDS"`
	REDIS_MAX_IDLE_CONN_NUMBER   int    `env:"REDIS_MAX_IDLE_CONN_NUMBER"`
	REDIS_MAX_ACTIVE_CONN_NUMBER int    `env:"REDIS_MAX_ACTIVE_CONN_NUMBER"`
}

type HttpConfig struct {
	Port int `env:"HTTP_PORT"`
}

type DBConfig struct {
	ConnectionURL          string `env:"DB_CONNECTION_URL"`
	Driver                 string `env:"DB_DRIVER"`
	MaxOpenConnNumber      int    `env:"DB_MAX_OPEN_CONN_NUMBER"`
	MaxIdleConnNumber      int    `env:"DB_MAX_IDLE_CONN_NUMBER"`
	ConnMaxLifeTimeSeconds int64  `env:"DB_CONN_MAX_LIFE_TIME_SECONDS"`
}
