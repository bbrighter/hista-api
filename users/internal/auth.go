package internal

import (
	"context"

	"encore.app/users/entity"
)

type (
	IAuthRepo interface {
		Login(ctx context.Context, name string, password string) (entity.User, []entity.UserAppPermission, error)
	}

	IAuth interface {
		Login(ctx context.Context, userName string, password string) (entity.User, []entity.UserAppPermission, error)
	}
)

type AuthUseCase struct {
	u IAuthRepo
}

func NewAuthUseCase(u IAuthRepo) AuthUseCase {
	return AuthUseCase{u: u}
}

func (uc AuthUseCase) Login(ctx context.Context, userName string, password string) (entity.User, []entity.UserAppPermission, error) {
	return uc.u.Login(ctx, userName, password)
}
