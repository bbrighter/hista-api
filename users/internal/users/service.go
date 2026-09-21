package users

import (
	"context"

	"encore.app/users/internal/shared"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type UserManager interface {
	CreateUser(ctx context.Context, user *shared.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, values map[string]any) error
	FindUser(ctx context.Context, id uuid.UUID) (shared.User, error)
	FindUserByName(ctx context.Context, name string) (shared.User, error)
	ListUsers(ctx context.Context) ([]shared.User, error)
	FindPermissions(ctx context.Context, id uuid.UUID) ([]shared.UserAppPermission, error)
}

type Encrypter interface {
	GeneratePassword(string) (string, error)
	ValidatePassword(string, string) error
}

type UserService struct {
	u UserManager
	e Encrypter
}

func NewUserService(db *gorm.DB, costs int) *UserService {
	r := NewUserRepo(db)
	e := NewEncryption(costs)
	return &UserService{u: r, e: e}
}

func (s *UserService) CreateUser(ctx context.Context, name string, password string) (shared.User, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return shared.User{}, err
	}
	hashedPassword, err := s.e.GeneratePassword(password)
	if err != nil {
		return shared.User{}, err
	}
	var user = shared.User{
		ID:       uuid,
		Name:     name,
		Password: hashedPassword,
	}
	if err = s.u.CreateUser(ctx, &user); err != nil {
		return shared.User{}, err
	}
	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string, oldPassword string) error {
	user, err := s.u.FindUser(ctx, id)
	if err != nil {
		return err
	}
	if err := s.e.ValidatePassword(user.Password, oldPassword); err != nil {
		return err
	}
	newPasswordHashed, err := s.e.GeneratePassword(newPassword)
	if err != nil {
		return err
	}
	return s.u.UpdateUser(ctx, id, map[string]any{"password": newPasswordHashed})
}

func (s *UserService) ChangePasswordForced(ctx context.Context, id uuid.UUID, newPassword string) error {
	_, err := s.u.FindUser(ctx, id)
	if err != nil {
		return err
	}
	newPasswordHashed, err := s.e.GeneratePassword(newPassword)
	if err != nil {
		return err
	}
	return s.u.UpdateUser(ctx, id, map[string]any{"password": newPasswordHashed})
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.u.DeleteUser(ctx, id)
}

func (s *UserService) Login(ctx context.Context, name string, password string) (shared.User, []shared.UserAppPermission, error) {
	user, err := s.u.FindUserByName(ctx, name)
	if err != nil {
		return shared.User{}, []shared.UserAppPermission{}, err
	}
	if err := s.e.ValidatePassword(user.Password, password); err != nil {
		return shared.User{}, []shared.UserAppPermission{}, err
	}
	permissions, err := s.u.FindPermissions(ctx, user.ID)
	if err != nil {
		return shared.User{}, []shared.UserAppPermission{}, err
	}
	return user, permissions, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]shared.User, error) {
	return s.u.ListUsers(ctx)
}

func (s *UserService) FindUserByName(ctx context.Context, name string) (shared.User, error) {
	return s.u.FindUserByName(ctx, name)
}

func (s *UserService) FindUser(ctx context.Context, id uuid.UUID) (shared.User, error) {
	return s.u.FindUser(ctx, id)
}

type SettingParams struct {
	LoadingMode *string
	Language    *string
}

func (s *UserService) ChangeSettings(ctx context.Context, id uuid.UUID, settings SettingParams) error {
	values := make(map[string]any)
	if settings.Language != nil {
		values["language"] = *settings.Language
	}
	if settings.LoadingMode != nil {
		values["loading_mode"] = *settings.LoadingMode
	}
	return s.u.UpdateUser(ctx, id, values)
}
