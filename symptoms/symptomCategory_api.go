package symptoms

import "context"

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
