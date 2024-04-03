package config

type KibanaConfig struct {
	URL string `yaml:"kibana_url" env-required:"true"`
}
