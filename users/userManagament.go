package users

import (
	"context"

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

// encore:api private method=POST path=/internal/user
func (s *Service) CreateUser(ctx context.Context, params UserParams) (UserIdResponse, error) {
	user, err := s.mgmt.Create(ctx, params.Name, params.Password)
	return UserIdResponse{UserId: user.ID}, err
}

// encore:api private method=DELETE path=/internal/user/:id
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.mgmt.Delete(ctx, id)
}

// encore:api private method=GET path=/internal/user
func (s *Service) ListUsers(ctx context.Context) (entity.UserListResponse, error) {
	users := s.mgmt.List(ctx)
	return users.ToResponse(), nil
}

// encore:api private method=GET path=/internal/user/:name
func (s *Service) Exists(ctx context.Context, name string) (UserIdResponse, error) {
	user, err := s.mgmt.Find(ctx, name)
	return UserIdResponse{UserId: user.ID}, err
}

type AddUserToProductInstanceParams struct {
	AppIds    []string `json:"appIds"`
	ProductId string   `json:"productId"`
}

// encore:api private method=POST path=/internal/user/:userId/product-instance/:productInstanceId
func (s *Service) AddUserToProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID, params AddUserToProductInstanceParams) error {
	return s.mgmt.AddUserToInstance(ctx, productInstanceId, params.ProductId, userId, params.AppIds)
}

// encore:api private method=DELETE path=/internal/user/:userId/product-instance/:productInstanceId
func (s *Service) RemoveUserFromProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID) error {
	return s.mgmt.RemoveUserFromInstance(ctx, productInstanceId, userId)
}

// encore:api private method=GET path=/internal/product-instance/:productInstanceId/users
func (s *Service) ListUsersForProductInstance(ctx context.Context, productInstanceId uuid.UUID) (entity.UserListResponse, error) {
	users, err := s.mgmt.ListForInstance(ctx, productInstanceId)
	return users.ToResponse(), err
}
