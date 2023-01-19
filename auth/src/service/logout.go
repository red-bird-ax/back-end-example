package service

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) logout(response http.ResponseWriter, request *http.Request) {
	requestBody, err := net.UnmarshalRequestBody[logoutRequestBody](request)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to parse request", http.StatusBadRequest)
		return
	}

	if err = requestBody.Validate(); err != nil {
		service.log.Error(err)
		net.ResponseError(response, err.Error(), http.StatusBadRequest)
		return
	}

	if err = service.cache.Del(service.ctx, requestBody.Token).Err(); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to delete refresh token", http.StatusInternalServerError)
	}
}

type logoutRequestBody struct {
	Token string `json:"refresh-token"`
}

func (request logoutRequestBody) Validate() error {
	return validation.Validate(&request.Token, validation.Required)
}