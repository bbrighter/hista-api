package users

import (
	"context"

	"encore.dev/types/uuid"
)

type UserPasswordChangeParams struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

// encore:api auth method=PATCH path=/user/:id tag:user
func (s *Service) PatchPassword(ctx context.Context, id uuid.UUID, params UserPasswordChangeParams) error {
	return s.user.ChangePassword(ctx, id, params.NewPassword, params.OldPassword)
}
