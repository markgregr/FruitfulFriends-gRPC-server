package config

import "time"

type JWTConfig struct {
	TokenTTL time.Duration `env:"GRPC_SERVER_TOKEN_TTL" env-required:"true"`
	TokenKey string        `env:"GRPC_SERVER_JWT_SECRET" env-required:"true"`
}
