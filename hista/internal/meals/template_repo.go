package meals

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type templateRepo struct {
	db *gorm.DB
}

func newTemplateRepo(db *gorm.DB) *templateRepo {
	return &templateRepo{db: db}
}

func (r *templateRepo) FirstTemplateAndItems(ctx context.Context, id uint) (Template, error) {
	return gorm.G[Template](r.db).
		Scopes(wherePiid(ctx)).
		Where("id = ?", id).
		Preload("Items", nil).First(ctx)
}

func (r *templateRepo) ListTemplatesAndItems(ctx context.Context) ([]Template, error) {
	return gorm.G[Template](r.db).
		Scopes(wherePiid(ctx)).
		Preload("Items", nil).Find(ctx)
}

func (r *templateRepo) DeleteTemplate(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Template](ctx, r.db, id)
}

func (r *templateRepo) CreateTemplate(ctx context.Context, template *Template) error {
	return generic_queries.Create(ctx, r.db, template)
}

func (r *templateRepo) ReplaceTemplate(ctx context.Context, id uint, name string, items []TemplateItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := generic_queries.UpdateColumn[*Template](ctx, tx, id, "name", name); err != nil {
			return err
		}
		if _, err := gorm.G[TemplateItem](tx).
			Where("template_id = ?", id).
			Delete(ctx); err != nil {
			return err
		}
		for i := range items {
			items[i].TemplateID = id
		}
		return gorm.G[TemplateItem](tx).CreateInBatches(ctx, &items, 100)
	})
}
