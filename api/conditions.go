package api

import (
	"context"

	entity "encore.app/entity"
)

type PatchSeverityRequestParams struct {
	Severity entity.Severity `json:"severity"`
}

// encore:api auth method=PATCH path=/conditions/:conditionID
func (service *Service) PatchCondition(ctx context.Context, conditionID uint, params PatchSeverityRequestParams) error {
	return service.conditions.PatchSeverity(ctx, conditionID, params.Severity)
}

// encore:api auth method=DELETE path=/conditions/:conditionID
func (service *Service) DeleteCondition(ctx context.Context, conditionID uint) (entity.SymptomCategoriesResponse, error) {
	cats, err := service.conditions.Delete(ctx, conditionID)
	return cats.ToResponse(), err

}
