package users

import (
	"context"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *User) error {
	return gorm.G[User](r.db).Create(ctx, user)
}

func (r *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	rows, err := gorm.G[User](r.db).Where("id = ?", id).Delete(ctx)
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

func (r UserRepo) FindUser(ctx context.Context, id uuid.UUID) (User, error) {
	return gorm.G[User](r.db).Where("id = ?", id).First(ctx)
}

func (r *UserRepo) FindPermissions(ctx context.Context, id uuid.UUID) ([]UserAppPermission, error) {
	return gorm.G[UserAppPermission](r.db).
		Preload("UserProductInstance", nil).
		Where("user_id = ?", id).
		Find(ctx)
}

func (r UserRepo) FindUserByName(ctx context.Context, name string) (User, error) {
	return gorm.G[User](r.db).Where("name = ?", name).First(ctx)
}

func (r UserRepo) ListUsers(ctx context.Context) ([]User, error) {
	return gorm.G[User](r.db).Find(ctx)

}

func (r *UserRepo) AddUserToProductInstance(ctx context.Context, instanceId uuid.UUID, productId string, user User, apps []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		upi := UserProductInstance{
			ProductId:         productId,
			ProductInstanceId: instanceId,
			UserId:            user.ID,
		}
		if err := gorm.G[UserProductInstance](tx).Create(ctx, &upi); err != nil {
			return err
		}

		var permissions []UserAppPermission
		for _, app := range apps {
			var perm = UserAppPermission{
				UserId:                user.ID,
				UserProductInstanceId: upi.ID,
				App:                   app,
				PermissionLevel:       "full",
			}
			permissions = append(permissions, perm)
		}
		return gorm.G[UserAppPermission](tx).CreateInBatches(ctx, &permissions, 100)
	})
}

func (r *UserRepo) RemoveUserFromProductInstance(ctx context.Context, instanceId uuid.UUID, user User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		upi, err := gorm.G[UserProductInstance](tx).Where("product_instance_id = ? AND user_id = ?", instanceId, user.ID).First(ctx)
		if err != nil {
			return err
		}
		_, err = gorm.G[UserAppPermission](tx).Where("user_product_instance_id = ? AND user_id = ?", upi.ID, user.ID).Delete(ctx)
		if err != nil {
			return err
		}
		_, err = gorm.G[UserProductInstance](tx).Where("id = ?", upi.ID).Delete(ctx)
		return err
	})
}

func (r *UserRepo) ListUserForInstance(ctx context.Context, instanceId uuid.UUID) ([]User, error) {
	usersForProdInstance, err := gorm.G[UserProductInstance](r.db).Where("product_instance_id = ?", instanceId).Find(ctx)
	if err != nil {
		return []User{}, err
	}
	if len(usersForProdInstance) == 0 {
		return []User{}, gorm.ErrRecordNotFound
	}
	var userIds []uuid.UUID
	for _, u := range usersForProdInstance {
		userIds = append(userIds, u.UserId)
	}
	return gorm.G[User](r.db).Where("id IN ?", userIds).Find(ctx)
}
