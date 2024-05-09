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
