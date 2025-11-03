package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvDevelopment = "dev"
	EnvTesting     = "test"
)

type Config struct {
	DBName      string `env:"DB_NAME"`
	DBUser      string `env:"DB_USER"`
	DBPassword  string `env:"DB_PASSWORD"`
	DBHost      string `env:"DB_HOST"`
	DBPort      string `env:"DB_PORT"`
	DBPortTest  string `env:"DB_PORT_TEST"`
	Env         Env    `env:"ENV" envDefault:"dev"`
	ProjectRoot string `env:"PROJECT_ROOT"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) DatabaseUrl() string {
	port := c.DBPort
	if c.Env == EnvTesting {
		port = c.DBPortTest
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		port,
		c.DBName,
	)
}
