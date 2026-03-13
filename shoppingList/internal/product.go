package internal

import (
	"context"
)

type IProductManager interface {
	Archive(ctx context.Context, id uint, archive bool) error
	PatchName(ctx context.Context, id uint, newName string) error
	Delete(ctx context.Context, id uint) error
}

type ProductManager struct {
	r ProductRepo
}

func NewProductManager(r ProductRepo) *ProductManager {
	return &ProductManager{r: r}
}

func (pm ProductManager) Archive(ctx context.Context, id uint, archive bool) error {
	values := map[string]any{
		"archived": archive,
	}
	return pm.r.Update(ctx, id, values)
}
func (pm ProductManager) PatchName(ctx context.Context, id uint, newName string) error {
	values := map[string]any{
		"name": newName,
	}
	return pm.r.Update(ctx, id, values)
}
func (pm ProductManager) Delete(ctx context.Context, id uint) error {
	return pm.r.Delete(ctx, id)
}
