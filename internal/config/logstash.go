package config

type LogstashConfig struct {
	Host string `yaml:"logstash_host" env-required:"true"`
	Port int    `yaml:"logstash_port" env-required:"true"`
}
