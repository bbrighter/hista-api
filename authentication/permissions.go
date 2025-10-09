package authentication

import (
	"context"

	"encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/auth"
)

// encore:api auth method=GET path=/permissions
func (s *Service) GetPermissions(ctx context.Context) (*entity.AuthData, error) {
	data, ok := auth.Data().(*entity.AuthData)
	if !ok {
		return data, errors.ErrorUnauthenticated
	}
	return data, nil
}
