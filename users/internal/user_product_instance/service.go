package userproductinstance

import (
	"context"

	"encore.app/users/internal/shared"
	"encore.dev/types/uuid"
)

type UserProductInstanceService struct {
	i *userProductInstanceRepo
}

func NewUserProductInstanceService(r *userProductInstanceRepo) *UserProductInstanceService {
	return &UserProductInstanceService{i: r}
}

func (s *UserProductInstanceService) AddUserToInstance(
	ctx context.Context,
	instanceId uuid.UUID,
	productId string,
	userId uuid.UUID,
	app_ids []string,
) error {
	return s.i.AddUserToProductInstance(ctx, instanceId, productId, userId, app_ids)
}

func (s *UserProductInstanceService) RemoveUserFromInstance(ctx context.Context, instanceId uuid.UUID, userId uuid.UUID) error {
	return s.i.RemoveUserFromProductInstance(ctx, instanceId, userId)
}

func (s *UserProductInstanceService) ListForInstance(ctx context.Context, instanceId uuid.UUID) (shared.Users, error) {
	return s.i.ListUserForInstance(ctx, instanceId)
}
