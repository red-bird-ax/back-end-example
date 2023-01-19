package service

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"

	"github.com/red-bird-ax/poster/users/src/repositroy/model"
	"github.com/red-bird-ax/poster/utils/jwt"
	net "github.com/red-bird-ax/poster/utils/network"
)

func (service *Service) unsubscribe(response http.ResponseWriter, request *http.Request) {
	id, _, err := jwt.GetUser(request)
	if err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to get user from jwt", http.StatusInternalServerError)
		return
	}

	requestBody, err := net.UnmarshalRequestBody[unsubscribeReqeust](request)
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

	if err = service.subs.Delete(model.Subscription{SubscriberID: id, UserID: requestBody.UserID}); err != nil {
		service.log.Error(err)
		net.ResponseError(response, "failed to unsubscribe user", http.StatusInternalServerError)
	}
}

type unsubscribeReqeust struct {
	UserID uuid.UUID `json:"user-id"`
}

func (request unsubscribeReqeust) Validate() error {
	return validation.Validate(&request.UserID, validation.Required)
}