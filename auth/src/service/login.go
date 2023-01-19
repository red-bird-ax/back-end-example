package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"

	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

var ErrUnauthorized = errors.New("authentication failed")

func (service *Service) login(response http.ResponseWriter, request *http.Request) {
	requestBody, err := net.UnmarshalRequestBody[loginRequestBody](request)
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

	id, name, err := service.authenticate(requestBody.Name, requestBody.Password)
	if err != nil {
		if err == ErrUnauthorized {
			service.log.Warning(err)
			response.WriteHeader(http.StatusUnauthorized)
		} else {
			service.log.Error(err)
			net.ResponseError(response, "failed to authenticate user", http.StatusInternalServerError)
		}
		return
	}

	accessToken, err := jwt.GenerateAccessToken(id, name, service.accessTokenOptions)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(id, name, service.refreshTokenOptions)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	if err = service.cache.Set(service.ctx, refreshToken, 0, service.refreshTokenOptions.Expire).Err(); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to save refresh token", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(response).Encode(loginResponseBody{AccessToken: accessToken, RefreshToken: refreshToken}); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to create login response", http.StatusInternalServerError)
	}
}

func (service *Service) authenticate(name, password string) (uuid.UUID, string, error) {
	type RequestBody struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	type ResponseBody struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	requestBody, err := json.Marshal(RequestBody{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return uuid.Nil, "", err
	}

	endpoint := "http://" + service.usersEndpoint + "/authenticate"
	response, err := service.client.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return uuid.Nil, "", err
	}

	if response.StatusCode != http.StatusOK {
		return uuid.Nil, "", ErrUnauthorized
	}

	if responseBody, err := net.UnmarshalResponseBody[ResponseBody](response); err == nil {
		return responseBody.ID, responseBody.Name, nil
	} else {
		return uuid.Nil, "", err
	}
}

type loginRequestBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (request loginRequestBody) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 255)),
	)
}

type loginResponseBody struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}