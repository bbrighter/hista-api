package api

import (
	"context"

	entity "encore.app/entity"
)

type PatchSeverityRequestParams struct {
	Severity entity.ConditionSeverity `json:"severity"`
}

// encore:api auth method=PATCH path=/conditions/:conditionID
func (service *Service) PatchCondition(ctx context.Context, conditionID uint, params PatchSeverityRequestParams) error {
	return service.conditions.PatchSeverity(conditionID, params.Severity)
}

// encore:api auth method=DELETE path=/conditions/:conditionID
func (service *Service) DeleteCondition(ctx context.Context, conditionID uint) (entity.SymptomCategoriesResponse, error) {
	cats, err := service.conditions.Delete(conditionID)
	return cats.ToResponse(), err

}
