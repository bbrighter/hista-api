package shoppinglist

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type ShoppingListRepo struct {
	db *gorm.DB
}

func NewShoppingListRepo(db *gorm.DB) *ShoppingListRepo {
	return &ShoppingListRepo{db: db}
}

func (r *ShoppingListRepo) CreateShoppingList(ctx context.Context, list *List) error {
	return generic_queries.Create(ctx, r.db, list)
}
func (r *ShoppingListRepo) DeleteShoppingList(ctx context.Context, id uint) error {
	return generic_queries.Delete[*List](ctx, r.db, id)
}
func (r *ShoppingListRepo) FirstShoppingList(ctx context.Context) (List, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return List{}, err
	}
	returnedList, err := gorm.G[List](r.db).
		Where("pi_id = ?", piid).
		Where("deleted_at IS NULL").
		Preload("Items", nil).
		First(ctx)
	return returnedList, err
}

func (r *ShoppingListRepo) CreateProduct(ctx context.Context, product *Product) error {
	return generic_queries.Create(ctx, r.db, product)
}
func (r *ShoppingListRepo) ListProducts(ctx context.Context) ([]*Product, error) {
	return generic_queries.List[*Product](ctx, r.db)
}
func (r *ShoppingListRepo) DeleteProduct(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Product](ctx, r.db, id)
}
func (r *ShoppingListRepo) UpdateProduct(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "products", id, values)
}

func (r *ShoppingListRepo) CreateItem(ctx context.Context, item *Item) error {
	return generic_queries.Create(ctx, r.db, item)
}
func (r *ShoppingListRepo) DeleteItems(ctx context.Context, ids []uint) error {
	return generic_queries.BatchDelete[*Item](ctx, r.db, ids)
}
func (r *ShoppingListRepo) UpdateItem(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "items", id, values)
}
func (r *ShoppingListRepo) ListItemsForList(ctx context.Context, listId uint) ([]Item, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return gorm.G[Item](r.db).
		Where("pi_id = ?", piid).
		Where("list_id = ?", listId).
		Find(ctx)
}
