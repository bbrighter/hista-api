package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/internal/symptoms"
	"encore.dev/types/uuid"
)

type PatchSeverityRequestParams struct {
	Severity uint `json:"severity"`
}

// encore:api auth method=PATCH path=/piid/:piid/conditions/:conditionID
func (service *Service) PatchCondition(ctx context.Context, piid uuid.UUID, conditionID uint, params PatchSeverityRequestParams) error {
	err := service.syms.PatchSeverity(ctx, conditionID, symptoms.Severity(params.Severity))
	return errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/conditions/:conditionID
func (service *Service) DeleteCondition(ctx context.Context, piid uuid.UUID, conditionID uint) (resp SymptomCategoryListResponse, err error) {
	cats, err := service.syms.DeleteCondition(ctx, conditionID)
	if err != nil {
		return resp, errors.MapError(err)
	}
	return toSymptomCategoryListResponse(cats), nil

}
