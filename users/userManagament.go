package users

import (
	"context"

	"encore.app/product_mgmt"
	"encore.app/users/entity"
	"encore.dev/types/uuid"
)

type UserParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserIdResponse struct {
	UserId uuid.UUID `json:"userId"`
}

// encore:api private method=POST path=/user
func (s *Service) CreateUser(ctx context.Context, params UserParams) (UserIdResponse, error) {
	user, err := s.mgmt.Create(ctx, params.Name, params.Password)
	return UserIdResponse{UserId: user.ID}, err
}

// encore:api private method=DELETE path=/user/:id
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.mgmt.Delete(ctx, id)
}

// encore:api private method=GET path=/user
func (s *Service) ListUsers(ctx context.Context) (entity.UserListResponse, error) {
	users := s.mgmt.List(ctx)
	return users.ToResponse(), nil
}

type AddUserToProductInstanceParams struct {
	AppIds []string `json:"appIds"`
}

// encore:api auth method=POST path=/user/:userId/product-instance/:productInstanceId tag:user
func (s *Service) AddUserToProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID, params AddUserToProductInstanceParams) error {
	resp, err := product_mgmt.FindInstance(ctx, productInstanceId)
	if err != nil {
		return err
	}
	return s.mgmt.AddUserToInstance(ctx, productInstanceId, resp.Product.ID, userId, params.AppIds)
}

// encore:api auth method=DELETE path=/user/:userId/product-instance/:productInstanceId tag:user
func (s *Service) RemoveUserFromProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID) error {
	_, err := product_mgmt.FindInstance(ctx, productInstanceId)
	if err != nil {
		return err
	}
	return s.mgmt.RemoveUserFromInstance(ctx, productInstanceId, userId)
}
