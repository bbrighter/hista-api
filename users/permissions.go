package users

import (
	"context"

	"encore.app/users/entity"
)

type LoginParams struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type PermissionResponse struct {
	User        entity.User                  `json:"user"`
	Permissions entity.UserAppPermissionList `json:"permissions"`
}

// encore:api private method=POST path=/internal/permissions
func (service *Service) GetPermissions(ctx context.Context, params LoginParams) (*PermissionResponse, error) {
	user, perm, err := service.auth.Login(ctx, params.UserName, params.Password)
	if err != nil {
		return &PermissionResponse{}, err
	}
	return &PermissionResponse{User: user, Permissions: perm}, nil
}
