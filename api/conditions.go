package api

import (
	"context"

	entity "encore.app/entity"
	"encore.dev/types/uuid"
)

type PatchSeverityRequestParams struct {
	Severity entity.Severity `json:"severity"`
}

// encore:api auth method=PATCH path=/piid/:piid/conditions/:conditionID tag:external
func (service *Service) PatchCondition(ctx context.Context, piid uuid.UUID, conditionID uint, params PatchSeverityRequestParams) error {
	return service.conditions.PatchSeverity(ctx, conditionID, params.Severity)
}

// encore:api auth method=DELETE path=/piid/:piid/conditions/:conditionID tag:external
func (service *Service) DeleteCondition(ctx context.Context, piid uuid.UUID, conditionID uint) (entity.SymptomCategoriesResponse, error) {
	cats, err := service.conditions.Delete(ctx, conditionID)
	return cats.ToResponse(), err

}
