package service

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Port               int           `env:"PORT"`
	UsersEndpoint      string        `env:"USERS_ENDPOINT"`
	RefreshTokenSecret string        `env:"REFRESH_TOKEN_SECRET"`
	RefreshTokenExpire time.Duration `env:"REFRESH_TOKEN_EXPIRE"`
	AccessTokenSecret  string        `env:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpire  time.Duration `env:"ACCESS_TOKEN_EXPIRE"`
	CachePassword      string        `env:"CACHE_PASSWORD"`
	CacheHost          string        `env:"CACHE_HOST"`
	CachePort          int           `env:"CACHE_PORT"`
	RequestTimeout     time.Duration `env:"REQUEST_TIMEOUT"`
}

func NewConfig() (*Config, error) {
	config := new(Config)
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}