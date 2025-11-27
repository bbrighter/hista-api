package repository

import (
	"context"

	"encore.app/users/entity"
	"encore.dev/types/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	db         *gorm.DB
	bcryptCost int
}

func NewUserRepo(db *gorm.DB, cost int) *UserRepo {
	return &UserRepo{db: db, bcryptCost: cost}
}

func (r UserRepo) Create(ctx context.Context, name string, password string) (entity.User, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return entity.User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), r.bcryptCost)
	if err != nil {
		return entity.User{}, err
	}
	var user = entity.User{
		ID:       uuid,
		Name:     name,
		Password: string(hashedPassword),
	}
	result := gorm.WithResult()
	err = gorm.G[entity.User](r.db, result).Create(ctx, &user)
	return user, err
}

func (r UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	rows, err := gorm.G[entity.User](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r UserRepo) ChangePassword(ctx context.Context, id uuid.UUID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), r.bcryptCost)
	if err != nil {
		return err
	}
	rows, err := gorm.G[entity.User](r.db).
		Where("id = ?", id).
		Update(ctx, "password", string(hashedPassword))
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r UserRepo) validatePassword(ctx context.Context, id uuid.UUID, password string) error {
	user, err := gorm.G[entity.User](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (r UserRepo) Login(ctx context.Context, name string, password string) (entity.User, []entity.UserAppPermission, error) {
	user, err := gorm.G[entity.User](r.db).
		Where("name = ?", name).
		First(ctx)
	if err != nil {
		return user, []entity.UserAppPermission{}, err
	}
	if err := r.validatePassword(ctx, user.ID, password); err != nil {
		return user, []entity.UserAppPermission{}, err
	}
	perm, err := gorm.G[entity.UserAppPermission](r.db).
		Preload("UserProductInstance", nil).
		Where("user_id = ?", user.ID).
		Find(ctx)
	return user, perm, err
}

func (r UserRepo) Find(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return gorm.G[entity.User](r.db).Where("id = ?", id).First(ctx)
}

func (r UserRepo) FindByName(ctx context.Context, name string) (entity.User, error) {
	return gorm.G[entity.User](r.db).Where("name = ?", name).First(ctx)
}

func (r UserRepo) List(ctx context.Context) entity.Users {
	users, _ := gorm.G[entity.User](r.db).Find(ctx)
	return entity.Users(users)
}

func (r *UserRepo) AddUser(ctx context.Context, instanceId uuid.UUID, productId string, user entity.User, apps []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		upi := entity.UserProductInstance{
			ProductId:         productId,
			ProductInstanceId: instanceId,
			UserId:            user.ID,
		}
		if err := gorm.G[entity.UserProductInstance](tx).Create(ctx, &upi); err != nil {
			return err
		}

		var permissions []entity.UserAppPermission
		for _, app := range apps {
			var perm = entity.UserAppPermission{
				UserId:                user.ID,
				UserProductInstanceId: upi.ID,
				App:                   app,
				PermissionLevel:       "full",
			}
			permissions = append(permissions, perm)
		}
		return gorm.G[entity.UserAppPermission](tx).CreateInBatches(ctx, &permissions, 100)
	})
}

func (r *UserRepo) RemoveUser(ctx context.Context, instanceId uuid.UUID, user entity.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		upi, err := gorm.G[entity.UserProductInstance](tx).Where("product_instance_id = ? AND user_id = ?", instanceId, user.ID).First(ctx)
		if err != nil {
			return err
		}
		_, err = gorm.G[entity.UserAppPermission](tx).Where("user_product_instance_id = ? AND user_id = ?", upi.ID, user.ID).Delete(ctx)
		if err != nil {
			return err
		}
		_, err = gorm.G[entity.UserProductInstance](tx).Where("id = ?", upi.ID).Delete(ctx)
		return err
	})
}

func (r *UserRepo) ListUserForInstance(ctx context.Context, instanceId uuid.UUID) ([]entity.User, error) {
	usersForProdInstance, err := gorm.G[entity.UserProductInstance](r.db).Where("product_instance_id = ?", instanceId).Find(ctx)
	if err != nil {
		return []entity.User{}, err
	}
	var userIds []uuid.UUID
	for _, u := range usersForProdInstance {
		userIds = append(userIds, u.UserId)
	}
	return gorm.G[entity.User](r.db).Where("id IN ?", userIds).Find(ctx)
}
