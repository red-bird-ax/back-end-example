package jwt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	jwtlib "github.com/golang-jwt/jwt"
	"github.com/red-bird-ax/poster/utils/network"
)

const (
	AccessTokenHeaderName  = "Access-Token"
	RefreshTokenHeaderName = "Refresh-Token"

	RefreshTokenEndpoint = "http://auth-service:1010/token"
)

type TokenOptions struct {
	Secret []byte
	Expire time.Duration
}

type claims struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	jwtlib.StandardClaims
}

func GenerateAccessToken(id uuid.UUID, name string, options TokenOptions) (string, error) {
	iat := time.Now().Unix()
	exp := time.Now().Add(options.Expire).Unix()

	accessToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims{
		ID:   id,
		Name: name,
		StandardClaims: jwtlib.StandardClaims{
			IssuedAt:  iat,
			ExpiresAt: exp,
		},
	})

	return accessToken.SignedString(options.Secret)
}

func GenerateRefreshToken(id uuid.UUID, name string, options TokenOptions) (string, error) {
	refreshToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims{
		ID:   id,
		Name: name,
	})
	return refreshToken.SignedString(options.Secret)
}

func Parse(tokenString string, options TokenOptions) (uuid.UUID, string, error) {
	token, err := jwtlib.Parse(tokenString, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return options.Secret, nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}

	if claims, ok := token.Claims.(jwtlib.MapClaims); ok && token.Valid {
		rawID := claims["id"].(string)
		name := claims["name"].(string)

		id, err := uuid.Parse(rawID)
		return id, name, err
	}
	return uuid.Nil, "", errors.New("invalid token")
}

func ParseFromRequest(request *http.Request, client *http.Client, options TokenOptions) (uuid.UUID, string, error) {
	accessToken := request.Header.Get(AccessTokenHeaderName)
	if accessToken == "" {
		return uuid.Nil, "", errors.New("access token not found in header")
	}

	id, name, err := Parse(accessToken, options)
	if err != nil {
		if vErr, ok := err.(*jwtlib.ValidationError); ok && (vErr.Errors & jwtlib.ValidationErrorExpired != 0) {
			refreshToken := request.Header.Get(RefreshTokenHeaderName)
			if refreshToken == "" {
				return uuid.Nil, "", errors.New("refresh token not found in header")
			}

			// todo: after access token is expired we refreshing it each request (beacuse we not return it)
			// need just to return err, user must fetch /token itself
			if accessToken, err = refreshAccessToken(refreshToken, client); err == nil {
				return Parse(accessToken, options)
			} else {
				return uuid.Nil, "", err
			}
		} else {
			return uuid.Nil, "", err
		}
	}
	return id, name, nil
}

func refreshAccessToken(refreshToken string, client *http.Client) (string, error) {
	type RequestBody struct {
		Token string `json:"refresh-token"`
	}

	type ResponseBody struct {
		Token string `json:"access-token"`
	}

	requestBody, err := json.Marshal(RequestBody{Token: refreshToken})
	if err != nil {
		return "", err
	}

	response, err := client.Post(RefreshTokenEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New("fail to refresh access token")
	}

	if responseBody, err := network.UnmarshalResponseBody[ResponseBody](response); err == nil {
		return responseBody.Token, nil
	} else {
		return "", err
	}
}

func SaveUser(request *http.Request, id uuid.UUID, name string) context.Context {
	ctx := context.WithValue(request.Context(), "id", id)
	ctx =  context.WithValue(ctx, "name", name)
	return ctx
}

func GetUser(request *http.Request) (uuid.UUID, string, error) {
	id, ok := request.Context().Value("id").(uuid.UUID)
	if !ok {
		return uuid.Nil, "", errors.New("failed to get user id from request")
	}
	name, ok := request.Context().Value("name").(string)
	if !ok {
		return uuid.Nil, "", errors.New("failed to get user name from request")
	}
	return id, name, nil
}