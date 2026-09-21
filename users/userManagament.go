package users

import (
	"context"

	"encore.app/errors"
	"encore.app/users/internal/shared"
	"encore.dev/types/uuid"
)

type UserParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserIdResponse struct {
	UserId uuid.UUID `json:"userId"`
}

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}

func toUserListResponse(us []shared.User) UserListResponse {
	var resp = []UserResponse{}
	for _, u := range us {
		resp = append(resp, UserResponse{
			ID:   u.ID,
			Name: u.Name,
		})
	}
	return UserListResponse{Users: resp}
}

// encore:api private method=POST path=/internal/user
func (s *Service) CreateUser(ctx context.Context, params UserParams) (UserIdResponse, error) {
	user, err := s.u.CreateUser(ctx, params.Name, params.Password)
	return UserIdResponse{UserId: user.ID}, errors.MapError(err)
}

// encore:api private method=DELETE path=/internal/user/:id
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.u.DeleteUser(ctx, id)
	return errors.MapError(err)
}

// encore:api private method=GET path=/internal/user
func (s *Service) ListUsers(ctx context.Context) (UserListResponse, error) {
	users, err := s.u.ListUsers(ctx)
	return toUserListResponse(users), errors.MapError(err)
}

// encore:api private method=GET path=/internal/user/:name
func (s *Service) Exists(ctx context.Context, name string) (UserIdResponse, error) {
	user, err := s.u.FindUserByName(ctx, name)
	return UserIdResponse{UserId: user.ID}, errors.MapError(err)
}

type AddUserToProductInstanceParams struct {
	AppIds    []string `json:"appIds"`
	ProductId string   `json:"productId"`
}

// encore:api private method=POST path=/internal/user/:userId/product-instance/:productInstanceId
func (s *Service) AddUserToProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID, params AddUserToProductInstanceParams) error {
	err := s.upi.AddUserToInstance(ctx, productInstanceId, params.ProductId, userId, params.AppIds)
	return errors.MapError(err)
}

// encore:api private method=DELETE path=/internal/user/:userId/product-instance/:productInstanceId
func (s *Service) RemoveUserFromProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID) error {
	err := s.upi.RemoveUserFromInstance(ctx, productInstanceId, userId)
	return errors.MapError(err)
}

// encore:api private method=GET path=/internal/product-instance/:productInstanceId/users
func (s *Service) ListUsersForProductInstance(ctx context.Context, productInstanceId uuid.UUID) (UserListResponse, error) {
	users, err := s.upi.ListForInstance(ctx, productInstanceId)
	return toUserListResponse(users), errors.MapError(err)
}
