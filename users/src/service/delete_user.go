package service

import (
	"net/http"

	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) deleteUser(response http.ResponseWriter, request *http.Request) { // todo: delete refresh token (via message queue)
	if id, _, err := jwt.GetUser(request); err == nil {
		if err := service.users.Delete(id); err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to delete user", http.StatusInternalServerError)
		}
	} else {
		service.log.Error(err)
		net.ResponseError(response, "failed to get user from jwt", http.StatusInternalServerError)
	}
}