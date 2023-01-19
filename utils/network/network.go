package network

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/red-bird-ax/poster/utils/data"
)

func ResponseError(response http.ResponseWriter, message string, status int) {
	type ResposeBody struct {
		Message string
		Error   string
	}
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(&ResposeBody{
		Message: message,
		Error:   http.StatusText(status),
	})
}

func UnmarshalRequestBody[T any](request *http.Request) (*T, error) {
	requestBody := new(T)
	if err := json.NewDecoder(request.Body).Decode(requestBody); err == nil {
		_ = request.Body.Close()
		return requestBody, nil
	} else {
		return nil, err
	}
}

func UnmarshalResponseBody[T any](response *http.Response) (*T, error) {
	responseBody := new(T)
	if err := json.NewDecoder(response.Body).Decode(responseBody); err == nil {
		_ = response.Body.Close()
		return responseBody, nil
	} else {
		return nil, err
	}
}

func GetIdFromURL(request *http.Request, name string) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(request, name))
}

func GetIntFromURL(request *http.Request, name string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(request, name), 10, 64)
}

func GetPaginationFromURL(request *http.Request) (*data.Pagination, error) {
	offset, err := strconv.ParseInt(request.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseInt(request.URL.Query().Get("amount"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &data.Pagination{
		Offset: offset,
		Limit:  amount,
	}, nil
}
