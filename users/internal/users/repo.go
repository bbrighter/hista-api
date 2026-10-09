package users

import (
	"context"

	"encore.app/users/internal/shared"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *shared.User) error {
	return gorm.G[shared.User](r.db).Create(ctx, user)
}

func (r *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	rows, err := gorm.G[shared.User](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, id uuid.UUID, values map[string]any) error {
	rows, err := gorm.G[map[string]any](r.db).Table("users").Where("id = ?", id).Updates(ctx, values)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r UserRepo) FindUser(ctx context.Context, id uuid.UUID) (shared.User, error) {
	return gorm.G[shared.User](r.db).Where("id = ?", id).First(ctx)
}

func (r *UserRepo) FindPermissions(ctx context.Context, id uuid.UUID) ([]shared.UserAppPermission, error) {
	return gorm.G[shared.UserAppPermission](r.db).
		Preload("UserProductInstance", nil).
		Where("user_id = ?", id).
		Find(ctx)
}

func (r UserRepo) FindUserByName(ctx context.Context, name string) (shared.User, error) {
	return gorm.G[shared.User](r.db).Where("name = ?", name).First(ctx)
}

func (r UserRepo) ListUsers(ctx context.Context) ([]shared.User, error) {
	return gorm.G[shared.User](r.db).Find(ctx)

}
