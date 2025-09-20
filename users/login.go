package users

import (
	"context"
	"time"

	"encore.app/authentication"
	"encore.app/errors"
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
	user, perm, err := service.auth.Login(ctx, params.UserName, params.Password)
	if err != nil {
		return &LoginResponse{}, errors.MapError(err)
	}
	var appIds []string
	var piid uuid.UUID
	for _, p := range perm {
		appIds = append(appIds, p.App)
		piid = p.UserProductInstance.ProductInstanceId
	}

	token := authentication.GenerateToken(user.Name, user.ID, appIds, piid, time.Hour*24)
	signedToken, err := token.SignedString(service.secrets.privateKey)
	if err != nil {
		return &LoginResponse{}, errors.MapError(err)
	}
	return &LoginResponse{Token: signedToken}, nil
}
