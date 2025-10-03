package internal

import (
	"context"

	"encore.app/users/entity"
	"encore.dev/types/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	IUserRepo interface {
		Find(ctx context.Context, id uuid.UUID) (entity.User, error)
		ChangePassword(ctx context.Context, id uuid.UUID, password string) error
	}

	IUser interface {
		ChangePassword(ctx context.Context, id uuid.UUID, newPassword string, oldPassword string) error
	}
)

type UserUseCase struct {
	r IUserRepo
}

func NewUserUseCase(r IUserRepo) UserUseCase {
	return UserUseCase{r: r}
}

func (uc UserUseCase) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string, oldPassword string) error {
	user, err := uc.r.Find(ctx, id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return err
	}
	return uc.r.ChangePassword(ctx, id, newPassword)
}
