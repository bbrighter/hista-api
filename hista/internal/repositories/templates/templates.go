package templates

import (
	"context"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

func (r *TemplateStore) List(ctx context.Context) (entity.Templates, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Templates{}, err
	}
	return gorm.G[*entity.Template](r.db).Where("pi_id = ?", piid).Preload("Items", nil).Find(ctx)
}
func (r *TemplateStore) Create(ctx context.Context, name string, items []entity.TemplateItem) (uint, error) {
	var template = entity.Template{Name: name, Items: items}
	err := generic_queries.Create(ctx, r.db, &template)
	return template.ID, err
}
func (r *TemplateStore) Update(ctx context.Context, id uint, name string, items []entity.TemplateItem) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	for _, i := range items {
		i.PIID = piid
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := generic_queries.UpdateColumn[*entity.Template](ctx, tx, id, "name", name); err != nil {
			return err
		}
		if _, err := gorm.G[entity.TemplateItem](tx).
			Where("pi_id = ?", piid).
			Where("template_id = ?", id).
			Delete(ctx); err != nil {
			return err
		}
		for i := range items {
			items[i].SetPiid(piid)
			items[i].TemplateID = id
		}
		return gorm.G[entity.TemplateItem](tx).CreateInBatches(ctx, &items, 100)
	})

}
func (r *TemplateStore) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Template](ctx, r.db, id)
}

func (r *TemplateStore) Get(ctx context.Context, id uint) (entity.Template, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Template{}, err
	}
	return gorm.G[entity.Template](r.db).
		Where("id = ?", id).
		Where("pi_id = ?", piid).
		Preload("Items", nil).First(ctx)

}
