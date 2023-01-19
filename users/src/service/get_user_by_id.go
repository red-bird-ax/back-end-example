package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) getUserByID(response http.ResponseWriter, request *http.Request) {
	id, err := net.GetIdFromURL(request, "id")
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to parse user id", http.StatusBadRequest)
		return
	}

	if user, err := service.users.Get(id); err == nil {
		responseBody := getUserResponse{
			ID:               user.ID,
			Name:             user.Name,
			FullName:         user.FullName,
			Status:           user.Status,
			RegistrationDate: user.RegitrationTimestamp,
		}

		if err := json.NewEncoder(response).Encode(&responseBody); err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to create user response", http.StatusInternalServerError)
		}
	} else {
		service.log.Error(err)
		net.ResponseError(response, "failed to get user", http.StatusInternalServerError)
	}
}

type getUserResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	FullName         string    `json:"full-name"`
	Status           string    `json:"status"`
	RegistrationDate time.Time `json:"registration-date"`
}