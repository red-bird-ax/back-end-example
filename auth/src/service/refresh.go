package service

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) refresh(response http.ResponseWriter, request *http.Request) {
	requestBody, err := net.UnmarshalRequestBody[refreshRequestBody](request)
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

	if service.cache.Get(service.ctx, requestBody.Token); err == nil {
		id, name, err := jwt.Parse(requestBody.Token, service.refreshTokenOptions)
		if err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to parse refresh token", http.StatusBadRequest)
			return
		}

		accessToken, err := jwt.GenerateAccessToken(id, name, service.accessTokenOptions)
		if err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		responseBody := refreshResponseBody{Token: accessToken}
		if err = json.NewEncoder(response).Encode(responseBody); err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to create refresh response", http.StatusInternalServerError)
		}
	} else {
		service.log.Warning(err)
		net.ResponseError(response, "token not found", http.StatusNotFound)
	}
}

type refreshRequestBody struct {
	Token string `json:"refresh-token"`
}

func (request refreshRequestBody) Validate() error {
	return validation.Validate(&request.Token, validation.Required)
}

type refreshResponseBody struct {
	Token string `json:"access-token"`
}