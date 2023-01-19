package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/red-bird-ax/poster/utils/data"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) searchUsers(response http.ResponseWriter, request *http.Request) {
	query := chi.URLParam(request, "query")
	if query == "" {
		service.log.Warningf("trying to find users with empty query")
		response.WriteHeader(http.StatusNoContent)
		return
	}

	var options *data.Options
	if pagination, err := net.GetPaginationFromURL(request); err == nil {
		options = &data.Options{
			Pagination: *pagination,
			OrderBy: "user_name",
		}
	}

	if users, err := service.users.SearchFor(query, options); err == nil {
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
		net.ResponseError(response, "failed to get users by quey", http.StatusInternalServerError)
	}
}