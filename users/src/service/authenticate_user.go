package service

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/red-bird-ax/poster/users/src/repositroy/model"
	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) authenticate(response http.ResponseWriter, request *http.Request) {
	requestBody, err := net.UnmarshalRequestBody[authenticateUserRequest](request)
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

	user, err := service.users.GetByName(requestBody.Name)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to get user", http.StatusInternalServerError)
		return
	}

	if isUserPassword(requestBody.Password, *user) {
		if err = json.NewEncoder(response).Encode(authenticateUserResponse{
			ID:   user.ID,
			Name: user.Name,
		}); err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to create user response", http.StatusInternalServerError)
			return
		}
	} else {
		response.WriteHeader(http.StatusUnauthorized)
	}
}

type authenticateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (request authenticateUserRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 255)),
	)
}

type authenticateUserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func isUserPassword(password string, user model.User) bool {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil
}

// JWT is implemented incorrectly!!!
func (service *Service) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if id, name, err := jwt.ParseFromRequest(request, &service.client, service.accessTokenOptions); err == nil {
			ctx := jwt.SaveUser(request, id, name)
			next.ServeHTTP(response, request.WithContext(ctx))
		} else {
			service.log.Error(err)
			net.ResponseError(response, "invalid or missing jwt", http.StatusUnauthorized)
		}
	})
}
