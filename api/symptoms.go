package api

import (
	"context"

	entity "encore.app/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/symptoms
func (service *Service) ListSymptoms(ctx context.Context, piid uuid.UUID) (entity.SymptomCategoriesResponse, error) {
	categories, err := service.symptoms.List(ctx)
	return categories.ToResponse(), err
}

type PatchSymptomNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/piid/:piid/symptoms/:id/name
func (service *Service) PatchSymptomName(ctx context.Context, piid uuid.UUID, id uint, params PatchSymptomNameParams) error {
	return service.symptoms.RenameSymptom(ctx, id, params.Name)
}

type PatchSymptomCategoryParams struct {
	ToCategoryID uint `json:"toCategoryId"`
}

// encore:api auth method=PATCH path=/piid/:piid/symptoms/:id/category
func (service *Service) PatchSymptomCategory(ctx context.Context, piid uuid.UUID, id uint, params PatchSymptomCategoryParams) error {
	return service.symptoms.ChangeCategory(ctx, id, params.ToCategoryID)
}

type PostSymptomCategoryRequest struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/piid/:piid/symptom-categories
func (service *Service) PostSymptomCategory(ctx context.Context, piid uuid.UUID, params PostSymptomCategoryRequest) (entity.IDResponse, error) {
	id, err := service.symptoms.CreateCategory(ctx, params.Name)
	return entity.IDResponse{ID: id}, err
}

type PatchCategoryNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/piid/:piid/symptom-categories/:id
func (service *Service) PatchCategoryName(ctx context.Context, piid uuid.UUID, id uint, params PatchCategoryNameParams) error {
	return service.symptoms.RenameCategory(ctx, id, params.Name)
}

// encore:api auth method=DELETE path=/piid/:piid/symptom-categories/:id
func (service *Service) DeleteSymptomCategory(ctx context.Context, piid uuid.UUID, id uint) error {
	return service.symptoms.DeleteCategory(ctx, id)
}
