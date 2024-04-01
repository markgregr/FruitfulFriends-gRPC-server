package config

import (
	"time"
)

type GRPCConfig struct {
	Port    int           `yaml: "port" env-required:"true"`
	Timeout time.Duration `yaml: "timeout" env-required:"true"`
}
