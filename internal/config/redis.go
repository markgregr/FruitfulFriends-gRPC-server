package config

type RedisConfig struct {
	Endpoint string `env:"GRPC_SERVER_REDIS_ENDPOINT" env-required:"true"`
	Password string `env:"GRPC_SERVER_REDIS_PASSWORD" env-required:"true"`
}
