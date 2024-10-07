package api

import (
	"context"

	entity "encore.app/entity"
)

// encore:api auth method=GET path=/symptoms
func (service Service) GetSymptoms(ctx context.Context) (entity.SymptomCategoriesResponse, error) {
	var categories entity.SymptomCategories = service.symptoms.List()
	return categories.ToResponse(), nil
}

type PostSymptomCategoryRequest struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/symptoms/categories
func (service *Service) PostSymptomCategory(ctx context.Context, params PostSymptomCategoryRequest) (entity.IDResponse, error) {
	id, err := service.symptoms.CreateCategory(params.Name)
	return entity.IDResponse{ID: id}, err
}
