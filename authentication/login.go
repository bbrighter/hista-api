package authentication

import (
	"context"

	"encore.app/users"
	"encore.dev/types/uuid"
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
		return &LoginResponse{}, nil
	}
	var appIds []string
	var piid uuid.UUID
	for _, perm := range resp.Permissions {
		appIds = append(appIds, perm.App)
		piid = perm.UserProductInstance.ProductInstanceId
	}
	signedToken, err := service.g.GenerateToken(params.UserName, resp.User.ID, appIds, piid)
	return &LoginResponse{Token: signedToken}, err
}
