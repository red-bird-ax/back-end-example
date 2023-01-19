package service

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"

	"github.com/red-bird-ax/poster/users/src/repositroy"
	"github.com/red-bird-ax/poster/utils/jwt"
	"github.com/red-bird-ax/poster/utils/logger"
)

type Service struct {
	router chi.Router
	users  repositroy.Users
	subs   repositroy.Subscriptions
	client http.Client
	log    logger.Logger
	port   int

	accessTokenOptions jwt.TokenOptions
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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", service.port), service.router); err != nil {
		fmt.Println("running auth users on port", service.port)
		service.log.Error(err)
	}
}

func (service *Service) setup(config Config) error {
	if err := service.setupDatabase(config); err != nil {
		return err
	}

	service.router = chi.NewRouter()
	service.log = logger.New()
	service.port = config.Port

	service.setupTokens(config)
	service.setupClient(config)
	service.setupRoutes()

	return nil
}

func (service *Service) setupRoutes()  {
	service.router.Post("/authenticate", service.authenticate)

	service.router.Post("/", service.createUser)
	service.router.Get("/{id}", service.getUserByID)
	service.router.Get("/", service.getAllUsers)
	service.router.Get("/search/{query}", service.searchUsers)
	service.router.With(service.authMiddleware).Patch("/", service.updateUser)
	service.router.With(service.authMiddleware).Delete("/", service.deleteUser)

	service.router.Route("/subscriptions", func(router chi.Router) {
		router.Get("/{id}", service.getUserSubscribtions) // who subscribed to this user
		router.With(service.authMiddleware).Post("/", service.subscribe)
		router.With(service.authMiddleware).Get("/", service.getSubscribtions) // who am i subscribed to
		router.With(service.authMiddleware).Delete("/", service.unsubscribe)
	})
}

func (service *Service) setupDatabase(config Config) error {
	connectionURL := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
	)

	if connection, err := dbx.MustOpen("postgres", connectionURL); err == nil {
		service.users = repositroy.NewUsers(connection)
		service.subs = repositroy.NewSubscriptions(connection)
		return nil
	} else {
		return err
	}
}

func (service *Service) setupClient(config Config) {
	service.client = http.Client{Timeout: config.RequestTimeout}
}

func (service *Service) setupTokens(config Config) {
	service.accessTokenOptions = jwt.TokenOptions{
		Secret: []byte(config.AccessTokenSecret),
	}
}