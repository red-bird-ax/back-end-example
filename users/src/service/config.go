package service

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Port               int           `env:"PORT"`
	DatabaseName       string        `env:"DB_NAME"`
	DatabasePort       string        `env:"DB_PORT"`
	DatabaseHost       string        `env:"DB_HOST"`
	DatabaseUser       string        `env:"DB_USER"`
	DatabasePassword   string        `env:"DB_PASSWORD"`
	AccessTokenSecret  string        `env:"ACCESS_TOKEN_SECRET"`
	RequestTimeout     time.Duration `env:"REQUEST_TIMEOUT"`
}

func NewConfig() (*Config, error) {
	config := new(Config)
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}