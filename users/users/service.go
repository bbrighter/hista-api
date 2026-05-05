package users

import (
	"context"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type UserManager interface {
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, values map[string]any) error
	FindUser(ctx context.Context, id uuid.UUID) (User, error)
	FindUserByName(ctx context.Context, name string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
}

type UserInstanceManager interface {
	AddUserToProductInstance(ctx context.Context, instanceId uuid.UUID, productId string, user User, apps []string) error
	RemoveUserFromProductInstance(ctx context.Context, instanceId uuid.UUID, user User) error
	ListUserForInstance(ctx context.Context, instanceId uuid.UUID) ([]User, error)
	FindPermissions(ctx context.Context, id uuid.UUID) ([]UserAppPermission, error)
}

type Encrypter interface {
	GeneratePassword(string) (string, error)
	ValidatePassword(string, string) error
}

type UserService struct {
	u UserManager
	i UserInstanceManager
	e Encrypter
}

func NewUserService(db *gorm.DB, costs int) *UserService {
	r := NewUserRepo(db)
	e := NewEncryption(costs)
	return &UserService{u: r, i: r, e: e}
}

func (s *UserService) CreateUser(ctx context.Context, name string, password string) (User, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}
	hashedPassword, err := s.e.GeneratePassword(password)
	if err != nil {
		return User{}, err
	}
	var user = User{
		ID:       uuid,
		Name:     name,
		Password: hashedPassword,
	}
	if err = s.u.CreateUser(ctx, &user); err != nil {
		return User{}, err
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

func (s *UserService) Login(ctx context.Context, name string, password string) (User, []UserAppPermission, error) {
	user, err := s.u.FindUserByName(ctx, name)
	if err != nil {
		return User{}, []UserAppPermission{}, err
	}
	if err := s.e.ValidatePassword(user.Password, password); err != nil {
		return User{}, []UserAppPermission{}, err
	}
	permissions, err := s.i.FindPermissions(ctx, user.ID)
	if err != nil {
		return User{}, []UserAppPermission{}, err
	}
	return user, permissions, nil
}

func (s *UserService) AddUserToInstance(ctx context.Context, instanceId uuid.UUID, productId string, userId uuid.UUID, app_ids []string) error {
	return s.i.AddUserToProductInstance(ctx, instanceId, productId, User{ID: userId}, app_ids)
}
func (s *UserService) RemoveUserFromInstance(ctx context.Context, instanceId uuid.UUID, userId uuid.UUID) error {
	return s.i.RemoveUserFromProductInstance(ctx, instanceId, User{ID: userId})
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.u.ListUsers(ctx)
}

func (s *UserService) FindUserByName(ctx context.Context, name string) (User, error) {
	return s.u.FindUserByName(ctx, name)
}

func (s *UserService) ListForInstance(ctx context.Context, instanceId uuid.UUID) (Users, error) {
	return s.i.ListUserForInstance(ctx, instanceId)
}
