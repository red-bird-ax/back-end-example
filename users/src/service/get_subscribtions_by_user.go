package service

import (
	"encoding/json"
	"net/http"

	"github.com/red-bird-ax/poster/utils/data"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) getUserSubscribtions(response http.ResponseWriter, request *http.Request) {
	id, err := net.GetIdFromURL(request, "id")
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to parse user id", http.StatusBadRequest)
		return
	}

	var options *data.Options
	if pagination, err := net.GetPaginationFromURL(request); err == nil {
		options = &data.Options{
			Pagination: *pagination,
		}
	}

	subscribtions, err := service.subs.GetByUser(id, options)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to get subscribtions", http.StatusInternalServerError)
		return
	}

	responseBody := make([]getSubscribtionsResponse, 0, len(subscribtions))
	for _, subscribtion := range subscribtions {
		responseBody = append(responseBody, getSubscribtionsResponse{UserID: subscribtion.SubscriberID})
	}

	if err = json.NewEncoder(response).Encode(responseBody); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to create subscribtions response", http.StatusInternalServerError)
	}
}