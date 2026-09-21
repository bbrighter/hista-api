package users

import (
	"context"

	"encore.app/errors"
	"encore.app/users/internal/shared"
	"encore.app/users/internal/users"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type UserSettingsResponse struct {
	LoadingMode string `json:"loadingMode"`
	Language    string `json:"language"`
}

func settingsToUserSettingsResponse(s shared.UserSettings) UserSettingsResponse {
	return UserSettingsResponse{
		LoadingMode: s.LoadingMode,
		Language:    s.Language,
	}
}

// encore:api auth method=GET path=/user-settings/:id
func (s *Service) GetUserSettings(ctx context.Context, id uuid.UUID) (UserSettingsResponse, error) {
	user, err := s.u.FindUser(ctx, id)
	return settingsToUserSettingsResponse(user.Settings), errors.MapError(err)
}

type UserSettingsPatchParams struct {
	LoadingMode option.Option[string] `json:"loadingMode"`
	Language    option.Option[string] `json:"language"`
}

// encore:api auth method=PATCH path=/user-settings/:id
func (s *Service) PatchUserSettings(ctx context.Context, id uuid.UUID, params UserSettingsPatchParams) error {
	err := s.u.ChangeSettings(ctx, id, users.SettingParams{
		LoadingMode: params.LoadingMode.PtrOrNil(),
		Language:    params.Language.PtrOrNil(),
	})
	return errors.MapError(err)
}
