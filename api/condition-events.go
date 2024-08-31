package api

import (
	"context"
	"time"

	entity "encore.app/entity"
)

type ConditionEventRequestParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=POST path=/condition-events
func (service *Service) CreateConditionEvent(ctx context.Context, params ConditionEventRequestParams) (entity.IDResponse, error) {
	id, err := service.conditionEvents.Create(params.Date)
	return entity.IDResponse{ID: id}, err
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

// encore:api auth method=PATCH path=/condition-events/:eventId
func (service *Service) PatchDate(ctx context.Context, eventId uint, params ConditionEventRequestParams) error {
	return service.conditionEvents.Patch(eventId, params.Date)
}

// encore:api auth method=DELETE path=/condition-events/:eventId
func (service *Service) DeleteConditionEvent(ctx context.Context, eventId uint) error {
	return service.conditionEvents.Delete(eventId)
}

type ConditionRequestParams struct {
	SymptomName *string `json:"symptomName" encore:"optional"`
	SymptomID   *uint   `json:"symptomId" encore:"optional"`
	CategoryID  *uint   `json:"categoryId" encore:"optional"`
}

// encore:api auth method=POST path=/condition-events/:eventId/conditions
func (service *Service) PostCondition(ctx context.Context, eventId uint, params ConditionRequestParams) (entity.PostConditionResponse, error) {
	var resp entity.PostConditionResponse
	condition, symptoms, err := service.conditions.Create(eventId, params.SymptomName, params.SymptomID, params.CategoryID)
	if err != nil {
		return resp, err
	}
	resp.Condition = condition.ToResponse()
	resp.Symptoms = symptoms.ToResponse()
	return resp, nil
}
