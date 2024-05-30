package symptoms

import (
	"context"

	"encore.app/errors"
)

type ConditionRequestParams struct {
	SymptomName *string `json:"symptomName" encore:"optional"`
	SymptomID   *uint   `json:"symptomId" encore:"optional"`
	CategoryID  uint    `json:"categoryId"`
}

type PostConditionResponse struct {
	Condition ConditionResponse         `json:"condition"`
	Symptoms  SymptomCategoriesResponse `json:"symptoms"`
}

// encore:api auth method=POST path=/condition-events/:eventId/conditions
func (service *Service) PostCondition(ctx context.Context, eventId uint, params ConditionRequestParams) (PostConditionResponse, error) {
	var resp PostConditionResponse
	if params.SymptomName == nil && params.SymptomID == nil {
		return resp, errors.ErrorAttributeMustBeSet("symptomName or symptomId")
	}
	if params.CategoryID == 0 {
		return resp, errors.ErrorAttributeMustBeSet("categoryId")
	}
	var err error
	var symptoms SymptomCategories
	var condition = Condition{ConditionEventID: eventId, Severity: Medium}
	if params.SymptomName != nil {
		symptoms, err = condition.createConditionBySymptomName(service, *params.SymptomName, params.CategoryID)
		if err != nil {
			return resp, err
		}
	} else if params.SymptomID != nil {
		condition.SymptomID = *params.SymptomID
		symptoms, err = condition.createConditionBySymptomID(service)
		if err != nil {
			return resp, err
		}
	}
	resp.Condition = condition.toResponse()
	resp.Symptoms = symptoms.toResponse()
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

type ConditionsResponse []ConditionResponse

// encore:api auth method=DELETE path=/conditions/:conditionID
func (service *Service) DeleteCondition(ctx context.Context, conditionID uint) (SymptomCategoriesResponse, error) {
	var condition = &Condition{ID: conditionID}
	var err error = condition.delete(service)
	var symptoms SymptomCategories = getSymptomCategories(service)
	return symptoms.toResponse(), err

}
