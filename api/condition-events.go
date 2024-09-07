package api

import (
	"context"
	"time"

	entity "encore.app/entity"
)

// encore:api auth method=POST path=/condition-events
func (service *Service) CreateConditionEvent(ctx context.Context) (entity.ConditionEventResponse, error) {
	event, err := service.conditionEvents.Create()
	return event.ToResponse(), err
}

// encore:api auth method=GET path=/condition-events
func (service *Service) GetConditionEvents(ctx context.Context) (entity.ConditionEventsResponse, error) {
	events := service.conditionEvents.List()
	return events.ToResponse(), nil
}

// encore:api auth method=GET path=/condition-events/:eventId
func (service *Service) GetConditionEvent(ctx context.Context, eventId uint) (entity.ConditionEventResponse, error) {
	event, err := service.conditionEvents.Get(eventId)
	return event.ToResponse(), err
}

type ConditionEventRequestParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=PATCH path=/condition-events/:eventId
func (service *Service) PatchDate(ctx context.Context, eventId uint, params ConditionEventRequestParams) error {
	return service.conditionEvents.Patch(eventId, params.Date)
}

// encore:api auth method=DELETE path=/condition-events/:eventId
func (service *Service) DeleteConditionEvent(ctx context.Context, eventId uint) (entity.SymptomCategoriesResponse, error) {
	cats, err := service.conditionEvents.Delete(eventId)
	return cats.ToResponse(), err
}

type ConditionRequestParams struct {
	SymptomName *string `json:"symptomName" encore:"optional"`
	SymptomID   *uint   `json:"symptomId" encore:"optional"`
	CategoryID  *uint   `json:"categoryId" encore:"optional"`
}

// encore:api auth method=POST path=/condition-events/:eventId/conditions
func (service *Service) PostCondition(ctx context.Context, eventId uint, params ConditionRequestParams) (resp entity.PostConditionResponse, err error) {
	condition, symptoms, err := service.conditions.Create(eventId, params.SymptomName, params.SymptomID, params.CategoryID)
	if err != nil {
		return resp, err
	}
	resp.Condition = condition.ToResponse()
	resp.Symptoms = symptoms.ToResponse()
	return resp, nil
}
