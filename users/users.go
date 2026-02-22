package users

import (
	"context"

	"encore.dev/types/uuid"
)

type UserPasswordChangeParams struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

// encore:api auth method=PATCH path=/user/:id
func (s *Service) PatchPassword(ctx context.Context, id uuid.UUID, params UserPasswordChangeParams) error {
	return s.user.ChangePassword(ctx, id, params.NewPassword, params.OldPassword)
}

type NewUserPasswordParams struct {
	NewPassword string `json:"newPassword"`
}

// encore:api private method=PATH path=/user/:id
func (s *Service) PatchPasswordWithoutValidation(ctx context.Context, id uuid.UUID, params NewUserPasswordParams) error {
	return s.user.ChangePasswordForced(ctx, id, params.NewPassword)
}
