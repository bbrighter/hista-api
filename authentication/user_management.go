package authentication

import (
	"context"

	"encore.app/product_mgmt"
	"encore.app/users"
	"encore.dev/types/uuid"
)

type UUIDResponse struct {
	ID uuid.UUID `json:"id"`
}

// encore:api private method=POST path=/product-instance/:productId/user/:userId
func (s Service) CreateProductInstanceWithOwner(ctx context.Context, productId string, userId uuid.UUID) (UUIDResponse, error) {
	if err := users.Exists(ctx, userId); err != nil {
		return UUIDResponse{}, err
	}
	if err := product_mgmt.FindProduct(ctx, productId); err != nil {
		return UUIDResponse{}, err
	}
	resp, err := product_mgmt.CreateInstance(ctx, product_mgmt.ProductInstanceParams{ProductId: productId, InstanceName: "name"})
	if err != nil {
		return UUIDResponse{}, err
	}
	if err := s.AddUserToProductInstance(ctx, userId, resp.ID); err != nil {
		return UUIDResponse{}, err
	}
	return UUIDResponse{ID: resp.ID}, nil

}

// encore:api auth method=POST path=/user/:userId/product-instance/:productInstanceId tag:user-management
func (s Service) AddUserToProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID) error {
	err := users.Exists(ctx, userId)
	if err != nil {
		return err
	}
	instance, err := product_mgmt.FindInstance(ctx, productInstanceId)
	if err != nil {
		return err
	}
	var appIds []string
	for _, app := range instance.Product.Apps {
		appIds = append(appIds, app.ID)
	}

	return users.AddUserToProductInstance(
		ctx,
		userId,
		productInstanceId,
		users.AddUserToProductInstanceParams{
			AppIds:    appIds,
			ProductId: instance.Product.ID,
		})
}

// encore:api auth method=DELETE path=/user/:userId/product-instance/:productInstanceId tag:user-management
func (s Service) RemoveUserFromProductInstance(ctx context.Context, userId uuid.UUID, productInstanceId uuid.UUID) error {
	if err := users.Exists(ctx, userId); err != nil {
		return err
	}
	_, err := product_mgmt.FindInstance(ctx, productInstanceId)
	if err != nil {
		return err
	}
	return users.RemoveUserFromProductInstance(ctx, userId, productInstanceId)
}
