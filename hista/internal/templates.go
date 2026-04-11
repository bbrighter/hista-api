package internal

import (
	"context"

	"encore.app/hista/entity"
)

type TemplateStore interface {
	Get(ctx context.Context, id uint) (entity.Template, error)
	List(ctx context.Context) (entity.Templates, error)
	Create(ctx context.Context, name string, items []entity.TemplateItem) (uint, error)
	Update(ctx context.Context, id uint, name string, items []entity.TemplateItem) error
	Delete(ctx context.Context, id uint) error
}

type FoodBatchStore interface {
	BatchCreateFoods(ctx context.Context, foods []entity.Food) ([]entity.Food, error)
}

type TemplateService interface {
	List(ctx context.Context) (entity.Templates, error)
	Create(ctx context.Context, name string, items []entity.TemplateItem) (uint, error)
	Update(ctx context.Context, id uint, name string, items []entity.TemplateItem) error
	Delete(ctx context.Context, id uint) error
	Apply(ctx context.Context, mealId uint, templateId uint) (entity.Foods, error)
}

type templateService struct {
	r TemplateStore
	f FoodBatchStore
}

func NewTemplateService(r TemplateStore, f FoodBatchStore) TemplateService {
	return templateService{r: r, f: f}
}

func (m templateService) List(ctx context.Context) (entity.Templates, error) {
	return m.r.List(ctx)
}
func (m templateService) Create(ctx context.Context, name string, items []entity.TemplateItem) (uint, error) {
	return m.r.Create(ctx, name, items)
}
func (m templateService) Update(ctx context.Context, id uint, name string, items []entity.TemplateItem) error {
	return m.r.Update(ctx, id, name, items)

}
func (m templateService) Delete(ctx context.Context, id uint) error {
	return m.r.Delete(ctx, id)
}

func (m templateService) Apply(ctx context.Context, mealId uint, templateId uint) (entity.Foods, error) {
	template, err := m.r.Get(ctx, templateId)
	if err != nil {
		return entity.Foods{}, err
	}
	var foods = []entity.Food{}
	for _, i := range template.Items {
		foods = append(foods, entity.Food{
			Condition:      i.Condition,
			IngredientID:   i.IngredientID,
			PIID:           i.PIID,
			IngredientPIID: i.PIID,
			MealPIID:       i.PIID,
			MealID:         mealId,
		})
	}

	foods, err = m.f.BatchCreateFoods(ctx, foods)

	var returnFoods = entity.Foods{}
	for _, f := range foods {
		returnFoods = append(returnFoods, &f)
	}
	return returnFoods, err
}
