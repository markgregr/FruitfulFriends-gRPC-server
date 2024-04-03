package config

import "time"

type JWTConfig struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	TokenKey string        `yaml:"token_key" env-required:"true"`
}
