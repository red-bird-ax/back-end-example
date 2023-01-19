package service

import (
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/red-bird-ax/poster/users/src/repositroy/model"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) createUser(response http.ResponseWriter, request *http.Request) {
	requestBody, err := net.UnmarshalRequestBody[createUserRequest](request)
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

	passwordHash, err := generatePasswordHash(requestBody.Password)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to hash password", http.StatusInternalServerError)
		return
	}

	userID := uuid.New()
	user := model.User{
		ID:                   userID,
		Name:                 requestBody.Name,
		FullName:             requestBody.FullName,
		PasswordHash:         passwordHash,
		Status:               "",
		RegitrationTimestamp: time.Now(),
	}

	if err = service.users.Create(user); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to save user", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(response).Encode(&createUserResponse{UserID: userID}); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to create user response", http.StatusInternalServerError)
	}
}

type createUserRequest struct {
	Name     string `json:"name"`
	FullName string `json:"full-name"`
	Password string `json:"password"`
}

func (request createUserRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&request.FullName, validation.Required, validation.Length(3, 255)),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 255)),
	)
}

type createUserResponse struct {
	UserID   uuid.UUID `json:"user-id"`
}

func generatePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}