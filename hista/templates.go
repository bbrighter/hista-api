package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/templates
func (service *Service) ListTemplates(ctx context.Context, piid uuid.UUID) (entity.TemplateListResponse, error) {
	templates, err := service.mealTemplates.List(ctx)
	return templates.ToResponse(), errors.MapError(err)
}

type TemplateParams struct {
	Name  string               `json:"name"`
	Items []TemplateItemParams `json:"items"`
}

type TemplateItemParams struct {
	IngredientId uint                 `json:"ingredientId"`
	Condition    entity.FoodCondition `json:"condition"`
}

func (p TemplateParams) ToItems() []entity.TemplateItem {
	var items = []entity.TemplateItem{}
	for _, i := range p.Items {
		items = append(items, entity.TemplateItem{
			IngredientID: i.IngredientId,
			Condition:    i.Condition,
		})
	}
	return items
}

// encore:api auth method=POST path=/piid/:piid/templates
func (service *Service) PostTemplate(ctx context.Context, piid uuid.UUID, params TemplateParams) (entity.IDResponse, error) {
	id, err := service.mealTemplates.Create(ctx, params.Name, params.ToItems())
	return entity.IDResponse{ID: id}, errors.MapError(err)
}

// encore:api auth method=PUT path=/piid/:piid/templates/:id
func (service *Service) PutTemplate(ctx context.Context, piid uuid.UUID, id uint, params TemplateParams) error {
	return errors.MapError(service.mealTemplates.Update(ctx, id, params.Name, params.ToItems()))
}

// encore:api auth method=DELETE path=/piid/:piid/templates/:id
func (service *Service) DeleteTemplate(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.mealTemplates.Delete(ctx, id))
}
