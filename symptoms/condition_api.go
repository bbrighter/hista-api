package symptoms

import (
	"context"

	"encore.app/errors"
)

type ConditionRequestParams struct {
	SymptomName string `json:"symptomName"`
	CategoryID  uint   `json:"categoryId"`
}

type PostConditionResponse struct {
	Condition ConditionResponse         `json:"condition"`
	Symptoms  SymptomCategoriesResponse `json:"symptoms"`
}

// encore:api auth method=POST path=/condition-events/:eventId/conditions
func (service *Service) PostCondition(ctx context.Context, eventId uint, params ConditionRequestParams) (PostConditionResponse, error) {
	var resp PostConditionResponse
	if params.SymptomName == "" {
		return resp, errors.ErrorAttributeMustBeSet("symptomName")
	}
	if params.CategoryID == 0 {
		return resp, errors.ErrorAttributeMustBeSet("categoryId")
	}
	var err error
	var symptoms SymptomCategories
	var condition = Condition{ConditionEventID: eventId, Severity: Medium}
	symptoms, err = condition.createConditionBySymptomName(service, params.SymptomName, params.CategoryID)
	resp.Condition = condition.toResponse()
	resp.Symptoms = symptoms.toResponse()
	return resp, err
}

// encore:api auth method=POST path=/condition-events/:eventId/conditions/symptoms/:symptomId
func (service *Service) PostConditionBySymptomID(ctx context.Context, eventId uint, symptomId uint) (ConditionResponse, error) {
	var condition = &Condition{SymptomID: symptomId, ConditionEventID: eventId}
	if err := condition.createConditionBySymptomID(service); err != nil {
		return ConditionResponse{}, err
	}
	var resp ConditionResponse = condition.toResponse()
	return resp, nil
}

type PatchSeverityRequestParams struct {
	Severity ConditionSeverity `json:"severity"`
}

// encore:api auth method=PATCH path=/conditions/:conditionID
func (service *Service) PatchCondition(ctx context.Context, conditionID uint, params PatchSeverityRequestParams) error {
	var condition = &Condition{ID: conditionID}
	return condition.changeSeverity(service, params.Severity)
}

// encore:api auth method=DELETE path=/conditions/:conditionID
func (service *Service) DeleteCondition(ctx context.Context, conditionID uint) error {
	var condition = &Condition{ID: conditionID}
	return condition.delete(service)
}
