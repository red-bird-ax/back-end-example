package service

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/red-bird-ax/poster/utils/data"
	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) getSubscribtions(response http.ResponseWriter, request *http.Request) {
	id, _, err := jwt.GetUser(request)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to get user from jwt", http.StatusInternalServerError)
		return
	}

	var options *data.Options
	if pagination, err := net.GetPaginationFromURL(request); err == nil {
		options = &data.Options{
			Pagination: *pagination,
		}
	}

	subscribtions, err := service.subs.GetBySubscriber(id, options)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to get subscribtions", http.StatusInternalServerError)
		return
	}

	responseBody := make([]getSubscribtionsResponse, 0, len(subscribtions))
	for _, subscribtion := range subscribtions {
		responseBody = append(responseBody, getSubscribtionsResponse{UserID: subscribtion.UserID})
	}

	if err = json.NewEncoder(response).Encode(responseBody); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to create subscribtions response", http.StatusInternalServerError)
	}
}

type getSubscribtionsResponse struct {
	UserID uuid.UUID `json:"user-id"`
}