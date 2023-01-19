package main

import (
	"log"

	"github.com/red-bird-ax/poster/auth/src/service"
)

func main()  {
	if authService, err := service.New(); err == nil {
		authService.Run()
	} else {
		log.Fatalf("failed to start auth service: %s", err.Error())
	}
}