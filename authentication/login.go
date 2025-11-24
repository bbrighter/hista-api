package authentication

import (
	"context"
	"time"

	"encore.app/authentication/entity"
	"encore.app/users"
	uuid "encore.dev/types/uuid"
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

	expirationTime := time.Hour * 24 * 7
	appMap := make(map[uuid.UUID]entity.ProductAndApps)
	for _, perm := range resp.Permissions {
		piid := perm.UserProductInstance.ProductInstanceId
		pa, ok := appMap[piid]
		if !ok {
			pa = entity.ProductAndApps{AppIds: []string{}, Product: perm.UserProductInstance.ProductId}
		}
		pa.AppIds = append(pa.AppIds, perm.App)
		appMap[piid] = pa
	}

	signedToken, err := service.g.GenerateToken(params.UserName, resp.User.ID, appMap, expirationTime)

	return &LoginResponse{Token: signedToken}, err
}
