package hista

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=POST path=/piid/:piid/condition-events tag:external
func (service *Service) CreateConditionEvent(ctx context.Context, piid uuid.UUID) (entity.ConditionEventResponse, error) {
	event, err := service.conditionEvents.Create(ctx)
	return event.ToResponse(), err
}

// encore:api auth method=GET path=/piid/:piid/condition-events tag:external
func (service *Service) ListConditionEvents(ctx context.Context, piid uuid.UUID) (entity.ConditionEventsResponse, error) {
	events, err := service.conditionEvents.List(ctx)
	return events.ToResponse(), err
}

// encore:api auth method=GET path=/piid/:piid/condition-events/:eventId tag:external
func (service *Service) GetConditionEvent(ctx context.Context, piid uuid.UUID, eventId uint) (entity.ConditionEventResponse, error) {
	event, err := service.conditionEvents.Get(ctx, eventId)
	return event.ToResponse(), err
}

type ConditionEventRequestParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=PATCH path=/piid/:piid/condition-events/:eventId tag:external
func (service *Service) PatchDate(ctx context.Context, piid uuid.UUID, eventId uint, params ConditionEventRequestParams) error {
	if params.Date.IsZero() {
		return errors.ErrorAttributeMustBeSet("date")
	}
	return service.conditionEvents.Patch(ctx, eventId, params.Date)
}

// encore:api auth method=DELETE path=/piid/:piid/condition-events/:eventId tag:external
func (service *Service) DeleteConditionEvent(ctx context.Context, piid uuid.UUID, eventId uint) (entity.SymptomCategoriesResponse, error) {
	cats, err := service.conditionEvents.Delete(ctx, eventId)
	return cats.ToResponse(), err
}

type ConditionRequestParams struct {
	SymptomName *string `json:"symptomName" encore:"optional"`
	SymptomID   *uint   `json:"symptomId" encore:"optional"`
	CategoryID  *uint   `json:"categoryId" encore:"optional"`
}

// encore:api auth method=POST path=/piid/:piid/condition-events/:eventId/conditions tag:external
func (service *Service) PostCondition(ctx context.Context, piid uuid.UUID, eventId uint, params ConditionRequestParams) (resp entity.PostConditionResponse, err error) {
	condition, symptoms, err := service.conditions.Create(ctx, eventId, params.SymptomName, params.SymptomID, params.CategoryID)
	if err != nil {
		return resp, err
	}
	resp.Condition = condition.ToResponse()
	resp.Symptoms = symptoms.ToResponse()
	return resp, nil
}
