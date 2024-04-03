package config

type ElasticSearchConfig struct {
	URL string `yaml:"elasticsearch_url" env-required:"true"`
}
