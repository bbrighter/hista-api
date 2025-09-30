package authentication

import (
	"context"

	"encore.app/users"
)

type LoginParams struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
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

	signedToken, err := service.g.GenerateToken(params.UserName, resp.User.ID, resp.Permissions.ToMap())
	return &LoginResponse{Token: signedToken}, err
}
