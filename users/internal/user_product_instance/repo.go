package userproductinstance

import (
	"context"

	"encore.app/users/internal/shared"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type userProductInstanceRepo struct {
	db *gorm.DB
}

func NewUserProductInstanceRepo(db *gorm.DB) *userProductInstanceRepo {
	return &userProductInstanceRepo{db: db}
}

func (r *userProductInstanceRepo) AddUserToProductInstance(ctx context.Context, instanceId uuid.UUID, productId string, userID uuid.UUID, apps []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		upi := shared.UserProductInstance{
			ProductId:         productId,
			ProductInstanceId: instanceId,
			UserId:            userID,
		}
		if err := gorm.G[shared.UserProductInstance](tx).Create(ctx, &upi); err != nil {
			return err
		}

		var permissions []shared.UserAppPermission
		for _, app := range apps {
			var perm = shared.UserAppPermission{
				UserId:                userID,
				UserProductInstanceId: upi.ID,
				App:                   app,
				PermissionLevel:       "full",
			}
			permissions = append(permissions, perm)
		}
		return gorm.G[shared.UserAppPermission](tx).CreateInBatches(ctx, &permissions, 100)
	})
}

func (r *userProductInstanceRepo) RemoveUserFromProductInstance(ctx context.Context, instanceId uuid.UUID, userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		upi, err := gorm.G[shared.UserProductInstance](tx).Where("product_instance_id = ? AND user_id = ?", instanceId, userID).First(ctx)
		if err != nil {
			return err
		}
		_, err = gorm.G[shared.UserAppPermission](tx).Where("user_product_instance_id = ? AND user_id = ?", upi.ID, userID).Delete(ctx)
		if err != nil {
			return err
		}
		_, err = gorm.G[shared.UserProductInstance](tx).Where("id = ?", upi.ID).Delete(ctx)
		return err
	})
}

func (r *userProductInstanceRepo) ListUserForInstance(ctx context.Context, instanceId uuid.UUID) ([]shared.User, error) {
	usersForProdInstance, err := gorm.G[shared.UserProductInstance](r.db).Where("product_instance_id = ?", instanceId).Find(ctx)
	if err != nil {
		return []shared.User{}, err
	}
	if len(usersForProdInstance) == 0 {
		return []shared.User{}, gorm.ErrRecordNotFound
	}
	var userIds []uuid.UUID
	for _, u := range usersForProdInstance {
		userIds = append(userIds, u.UserId)
	}
	return gorm.G[shared.User](r.db).Where("id IN ?", userIds).Find(ctx)
}
