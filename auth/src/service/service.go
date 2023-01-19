package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v9"

	"github.com/red-bird-ax/poster/utils/jwt"
	"github.com/red-bird-ax/poster/utils/logger"
)

type Service struct {
	router chi.Router
	log    logger.Logger
	client http.Client
	cache  *redis.Client
	port   int
	ctx    context.Context


	accessTokenOptions  jwt.TokenOptions
	refreshTokenOptions jwt.TokenOptions

	usersEndpoint string
}

func New() (*Service, error) {
	if config, err := NewConfig(); err == nil {
		var service Service
		if err = service.setup(*config); err == nil {
			return &service, err
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (service *Service) Run() {
	fmt.Println("running auth service on port", service.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", service.port), service.router); err != nil {
		service.log.Error(err)
	}
}

func (service *Service) setup(config Config) error {
	service.port = config.Port
	service.usersEndpoint = config.UsersEndpoint
	service.router = chi.NewRouter()
	service.log = logger.New()

	service.setupTokens(config)
	service.setupClient(config)
	service.setupRoutes()
	return service.setupCache(config)
}

func (service *Service) setupRoutes() {
	service.router.Delete("/logout", service.logout)
	service.router.Post("/login", service.login)
	service.router.Post("/token", service.refresh)
}

func (service *Service) setupClient(config Config) {
	service.client = http.Client{Timeout: config.RequestTimeout}
}

func (service *Service) setupTokens(config Config) {
	service.accessTokenOptions = jwt.TokenOptions{
		Secret: []byte(config.AccessTokenSecret),
		Expire: config.AccessTokenExpire,
	}
	service.refreshTokenOptions = jwt.TokenOptions{
		Secret: []byte(config.RefreshTokenSecret),
		Expire: config.RefreshTokenExpire,
	}
}

func (service *Service) setupCache(config Config) error {
	service.ctx = context.Background()
	time.Sleep(time.Second * 10) // todo: do actual healthcheck
	service.cache = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.CacheHost, config.CachePort),
		Password: config.CachePassword,
	})

	if pong, err := service.cache.Ping(service.ctx).Result(); err != nil {
		return err
	} else if pong != "PONG" {
		return errors.New("fail to ping redis")
	}

	return nil
}