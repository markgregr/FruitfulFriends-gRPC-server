package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Env      string `env:"GRPC_SERVER_ENV" env-default:"local"`
	LogsPath string `env:"GRPC_SERVER_LOGS_PATH_FILE" env-required:"true"`
	GRPC     GRPCConfig
	Postgres PostgresConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

func MustLoad() *Config {
	var cfg Config
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("unable to load .env file: %v", err)
		}
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error parsing environment variables: %v", err)
	}
	return &cfg
}
