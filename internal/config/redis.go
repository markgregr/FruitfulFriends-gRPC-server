package config

type RedisConfig struct {
	Endpoint string `yaml:"redis_endpoint" env-required:"true"`
	Password string `yaml:"redis_password" env-required:"true"`
}
