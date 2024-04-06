package config

import "time"

type GRPCConfig struct {
	Host    string        `env:"GRPC_SERVER_HOST" env-required:"true"`
	Port    int           `env:"GRPC_SERVER_PORT" env-required:"true"`
	Timeout time.Duration `env:"GRPC_SERVER_TIMEOUT" env-required:"true"`
}
