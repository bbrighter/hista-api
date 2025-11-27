package internal

import (
	"context"

	"encore.app/errors"
	"encore.app/users/entity"
	"encore.dev/types/uuid"
)

type (
	IUserManagementRepo interface {
		Create(ctx context.Context, name string, password string) (entity.User, error)
		Delete(ctx context.Context, id uuid.UUID) error
		AddUser(ctx context.Context, instanceId uuid.UUID, productId string, user entity.User, app_ids []string) error
		RemoveUser(ctx context.Context, instanceId uuid.UUID, user entity.User) error
		List(ctx context.Context) entity.Users
		FindByName(ctx context.Context, name string) (entity.User, error)
		ListUserForInstance(ctx context.Context, instanceId uuid.UUID) ([]entity.User, error)
	}

	IUserManagement interface {
		Create(ctx context.Context, name string, password string) (entity.User, error)
		Delete(ctx context.Context, id uuid.UUID) error
		AddUserToInstance(ctx context.Context, instanceId uuid.UUID, productId string, userId uuid.UUID, app_ids []string) error
		RemoveUserFromInstance(ctx context.Context, instanceId uuid.UUID, userId uuid.UUID) error
		List(ctx context.Context) entity.Users
		Find(ctx context.Context, name string) (entity.User, error)
		ListForInstance(ctx context.Context, instanceId uuid.UUID) (entity.Users, error)
	}
)

type UserManagementUseCase struct {
	r IUserManagementRepo
}

func NewUserManagement(r IUserManagementRepo) UserManagementUseCase {
	return UserManagementUseCase{r: r}
}

func (uc UserManagementUseCase) Create(ctx context.Context, name string, password string) (entity.User, error) {
	return uc.r.Create(ctx, name, password)
}

func (uc UserManagementUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.r.Delete(ctx, id)
}

func (uc UserManagementUseCase) AddUserToInstance(ctx context.Context, instanceId uuid.UUID, productId string, userId uuid.UUID, app_ids []string) error {
	return uc.r.AddUser(ctx, instanceId, productId, entity.User{ID: userId}, app_ids)
}
func (uc UserManagementUseCase) RemoveUserFromInstance(ctx context.Context, instanceId uuid.UUID, userId uuid.UUID) error {
	return uc.r.RemoveUser(ctx, instanceId, entity.User{ID: userId})
}

func (uc UserManagementUseCase) List(ctx context.Context) entity.Users {
	return uc.r.List(ctx)
}

func (uc UserManagementUseCase) Find(ctx context.Context, name string) (entity.User, error) {
	user, err := uc.r.FindByName(ctx, name)
	return user, errors.MapError(err)
}

func (uc UserManagementUseCase) ListForInstance(ctx context.Context, instanceId uuid.UUID) (entity.Users, error) {
	users, err := uc.r.ListUserForInstance(ctx, instanceId)
	return users, errors.MapError(err)
}
