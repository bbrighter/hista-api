package authentication

import (
	"context"
	"net/http"
	"time"

	"encore.app/users"
)

type LoginParams struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Cookie string `header:"Set-Cookie"`
}

// encore:api public method=POST path=/login
func (service *Service) Login(ctx context.Context, params LoginParams) (*LoginResponse, error) {
	resp, err := users.GetPermissions(ctx, users.LoginParams{
		UserName: params.UserName,
		Password: params.Password,
	})
	if err != nil {
		return &LoginResponse{}, err
	}

	expirationTime := time.Hour * 24 * 7
	signedToken, err := service.g.GenerateToken(params.UserName, resp.User.ID, resp.Permissions.ToMap(), expirationTime)

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    signedToken,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(expirationTime),
		HttpOnly: true,
	}

	return &LoginResponse{Cookie: cookie.String()}, err
}
