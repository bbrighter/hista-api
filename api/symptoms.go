package api

import (
	"context"

	entity "encore.app/entity"
)

// encore:api auth method=GET path=/symptoms
func (service *Service) GetSymptoms(ctx context.Context) (entity.SymptomCategoriesResponse, error) {
	categories, err := service.symptoms.List(ctx)
	return categories.ToResponse(), err
}

type PatchSymptomNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/symptoms/:id/name
func (service *Service) PatchSymptomName(ctx context.Context, id uint, params PatchSymptomNameParams) error {
	return service.symptoms.RenameSymptom(ctx, id, params.Name)
}

type PatchSymptomCategoryParams struct {
	ToCategoryID uint `json:"toCategoryId"`
}

// encore:api auth method=PATCH path=/symptoms/:id/category
func (service *Service) PatchSymptomCategory(ctx context.Context, id uint, params PatchSymptomCategoryParams) error {
	return service.symptoms.ChangeCategory(ctx, id, params.ToCategoryID)
}

type PostSymptomCategoryRequest struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/symptom-categories
func (service *Service) PostSymptomCategory(ctx context.Context, params PostSymptomCategoryRequest) (entity.IDResponse, error) {
	id, err := service.symptoms.CreateCategory(ctx, params.Name)
	return entity.IDResponse{ID: id}, err
}

type PatchCategoryNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/symptom-categories/:id
func (service *Service) PatchCategoryName(ctx context.Context, id uint, params PatchCategoryNameParams) error {
	return service.symptoms.RenameCategory(ctx, id, params.Name)
}

// encore:api auth method=DELETE path=/symptom-categories/:id
func (service *Service) DeleteSymptomCategory(ctx context.Context, id uint) error {
	return service.symptoms.DeleteCategory(ctx, id)
}
