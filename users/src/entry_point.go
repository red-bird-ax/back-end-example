package main

import (
	"log"

	"github.com/red-bird-ax/poster/users/src/service"
)

func main()  {
	if usersService, err := service.New(); err == nil {
		usersService.Run()
	} else {
		log.Fatalf("failed to start users service: %s", err.Error())
	}
}