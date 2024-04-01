package config

type PostgresConfig struct {
	URL            string `yaml:"postgres_url" env-required:"true"`
	AutoMigrate    bool   `yaml:"postgres_auto_migrate" env-default:"false"`
	MigrationsPath string `yaml:"postgres_migrations_path" env-required:"true"`
}
