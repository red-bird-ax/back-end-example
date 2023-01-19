package service

import (
	"encoding/json"
	"net/http"

	"github.com/red-bird-ax/poster/utils/data"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) getAllUsers(response http.ResponseWriter, request *http.Request) {
	var options *data.Options
	if pagination, err := net.GetPaginationFromURL(request); err == nil {
		options = &data.Options{
			Pagination: *pagination,
		}
	}

	if users, err := service.users.GetAll(options); err == nil {
		responseBody := make([]getUserResponse, 0, len(users))
		for _, user := range users {
			responseBody = append(responseBody, getUserResponse{
				ID:               user.ID,
				Name:             user.Name,
				FullName:         user.FullName,
				Status:           user.Status,
				RegistrationDate: user.RegitrationTimestamp,
			})
		}

		if err := json.NewEncoder(response).Encode(responseBody); err != nil {
			service.log.Error(err)
			net.ResponseError(response, "failed to create users response", http.StatusInternalServerError)
		}
	} else {
		service.log.Error(err)
		net.ResponseError(response, "failed to get users", http.StatusInternalServerError)
	}
}