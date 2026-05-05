package users

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

type UserPasswordChangeParams struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

// encore:api auth method=PATCH path=/user/:id
func (s *Service) PatchPassword(ctx context.Context, id uuid.UUID, params UserPasswordChangeParams) error {
	err := s.u.ChangePassword(ctx, id, params.NewPassword, params.OldPassword)
	return errors.MapError(err)
}

type NewUserPasswordParams struct {
	NewPassword string `json:"newPassword"`
}

// encore:api private method=PATH path=/user/:id
func (s *Service) PatchPasswordWithoutValidation(ctx context.Context, id uuid.UUID, params NewUserPasswordParams) error {
	err := s.u.ChangePasswordForced(ctx, id, params.NewPassword)
	return errors.MapError(err)
}
