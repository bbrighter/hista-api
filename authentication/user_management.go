package authentication

import (
	"context"

	"encore.app/errors"
	"encore.app/product_mgmt"
	entity "encore.app/shared/entity"
	"encore.app/users"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

type UUIDResponse struct {
	ID uuid.UUID `json:"id"`
}

// encore:api private method=POST path=/internal/product/:productId/users/:name
func (s Service) CreateProductInstanceWithOwner(ctx context.Context, productId string, name string) (UUIDResponse, error) {
	_, err := users.Exists(ctx, name)
	if err != nil {
		return UUIDResponse{}, err
	}
	if _, err := product_mgmt.FindProduct(ctx, productId); err != nil {
		return UUIDResponse{}, err
	}
	resp, err := product_mgmt.CreateInstance(ctx, product_mgmt.ProductInstanceParams{ProductId: productId, InstanceName: "name"})
	if err != nil {
		return UUIDResponse{}, err
	}
	return s.AddUserToProductInstance(ctx, resp.ID, name)
}

// encore:api auth method=POST path=/piid/:productInstanceId/users/:name
func (s Service) AddUserToProductInstance(ctx context.Context, productInstanceId uuid.UUID, name string) (UUIDResponse, error) {
	userId, err := users.Exists(ctx, name)
	if err != nil {
		return UUIDResponse{}, err
	}
	instance, err := product_mgmt.FindInstance(ctx, productInstanceId)
	if err != nil {
		return UUIDResponse{}, err
	}
	var appIds []string
	for _, app := range instance.Product.Apps {
		appIds = append(appIds, app.ID)
	}

	err = users.AddUserToProductInstance(
		ctx,
		userId.UserId,
		productInstanceId,
		users.AddUserToProductInstanceParams{
			AppIds:    appIds,
			ProductId: instance.Product.ID,
		})

	return UUIDResponse{ID: userId.UserId}, err
}

// encore:api auth method=DELETE path=/piid/:productInstanceId/users/:name
func (s Service) RemoveUserFromProductInstance(ctx context.Context, productInstanceId uuid.UUID, name string) error {
	data, ok := auth.Data().(*entity.AuthData)
	if !ok {
		return errors.NewEncoreError("no user", errs.InvalidArgument)
	}
	if data.UserName == name {
		return errors.NewEncoreError("cannot delete current user", errs.InvalidArgument)
	}
	userId, err := users.Exists(ctx, name)
	if err != nil {
		return err
	}
	_, err = product_mgmt.FindInstance(ctx, productInstanceId)
	if err != nil {
		return err
	}
	return users.RemoveUserFromProductInstance(ctx, userId.UserId, productInstanceId)
}

type UserResponse struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}

// encore:api auth method=GET path=/piid/:productInstanceId/users
func (s Service) GetUsersForProductInstance(ctx context.Context, productInstanceId uuid.UUID) (UserListResponse, error) {
	users, err := users.ListUsersForProductInstance(ctx, productInstanceId)
	var resp []UserResponse
	for _, u := range users.Users {
		resp = append(resp, UserResponse{Name: u.Name, ID: u.ID})
	}
	return UserListResponse{Users: resp}, err
}
