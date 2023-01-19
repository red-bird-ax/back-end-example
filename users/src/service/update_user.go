package service

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/red-bird-ax/poster/users/src/repositroy/model"
	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) updateUser(response http.ResponseWriter, request *http.Request) {
	id, _, err := jwt.GetUser(request)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to get user from jwt", http.StatusInternalServerError)
		return
	}

	requestBody, err := net.UnmarshalRequestBody[updateUserRequest](request)
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

	user, err := service.users.Get(id)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to parse request", http.StatusBadRequest)
		return
	}

	if err = service.users.Update(createUpdatedUser(*user, *requestBody)); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to update user", http.StatusInternalServerError)
	}
}

type updateUserRequest struct {
	Name     *string `json:"name"`
	FullName *string `json:"full-name"`
	Status   *string `json:"status"`
}

func (request updateUserRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Length(3, 100)),
		validation.Field(&request.FullName, validation.Length(3, 255)),
		validation.Field(&request.Status, validation.Length(0, 255)),
	)
}

func createUpdatedUser(user model.User, update updateUserRequest) model.User {
	if update.Name != nil {
		user.Name = *update.Name
	}
	if update.FullName != nil {
		user.FullName = *update.FullName
	}
	if update.Status != nil {
		user.Status = *update.Status
	}
	return user
}