package symptoms

import "context"

type SymptomResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CategoryID uint   `json:"categoryId"`
}

type SymptomCategoryResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Symptoms []SymptomResponse `json:"symptoms"`
}

type SymptomCategoriesResponse struct {
	Categories []SymptomCategoryResponse
}

// encore:api auth method=GET path=/symptoms
func (service Service) GetSymptoms(ctx context.Context) (SymptomCategoriesResponse, error) {
	var categories SymptomCategories = getSymptomCategories(&service)
	return categories.toResponse(), nil
}

type PostSymptomCategoryRequest struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/symptoms/categories
func (service *Service) PostSymptomCategory(ctx context.Context, params PostSymptomCategoryRequest) (IDResponse, error) {
	var category = &SymptomCategory{
		Name: params.Name,
	}
	var err error = category.create(service)
	return IDResponse{ID: category.ID}, err
}
