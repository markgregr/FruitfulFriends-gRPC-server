package config

type PostgresConfig struct {
	URL         string `env:"GRPC_SERVER_POSTGRES_URL" env-required:"true"`
	AutoMigrate bool   `env:"GRPC_SERVER_POSTGRES_AUTO_MIGRATE" env-default:"false"`
}
