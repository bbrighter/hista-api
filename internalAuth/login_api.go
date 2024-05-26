package internalAuth

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

type LoginParams struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token uuid.UUID `json:"token"`
}

//encore:api public method=POST path=/login
func (service *Service) Login(ctx context.Context, params LoginParams) (LoginResponse, error) {

	var user *User
	var err error
	user, err = getUserByName(params.UserName)
	if err != nil {
		return LoginResponse{}, err
	}
	var valid bool
	valid, err = user.isValidPassword(params.Password)
	if !valid {
		return LoginResponse{}, errors.ErrorUnauthenticated
	}
	var token Token
	token, err = user.firstOrCreateValidToken(service)
	return LoginResponse{Token: token.Bearer}, err
}
