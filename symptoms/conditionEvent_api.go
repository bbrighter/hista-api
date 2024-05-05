package symptoms

import (
	"context"
	"time"

	"encore.dev/beta/errs"
)

type ConditionEventRequestParams struct {
	Date time.Time `json:"date"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

// encore:api auth method=POST path=/condition-events
func (service *Service) CreateConditionEvent(ctx context.Context, params ConditionEventRequestParams) (IDResponse, error) {
	var event *ConditionEvent = newConditionEvent(params.Date)
	var err error = event.create(service)
	return IDResponse{ID: event.ID}, err
}

type ConditionEventsResponse struct {
	ConditionEvents []ConditionEventMetaResponse `json:"conditionEvents"`
}

type ConditionEventMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

// encore:api auth method=GET path=/condition-events
func (service *Service) GetConditionEvents(ctx context.Context) (ConditionEventsResponse, error) {
	return getConditionEvents(service).toResponse(), nil
}

type ConditionEventResponse struct {
	ID         uint                `json:"id"`
	Date       time.Time           `json:"date"`
	Conditions []ConditionResponse `json:"conditions"`
}

type ConditionResponse struct {
	ID       uint              `json:"id"`
	Symptom  SymptomResponse   `json:"symptom"`
	Severity ConditionSeverity `json:"severity"`
}

// encore:api auth method=GET path=/condition-events/:eventId
func (service *Service) GetConditionEvent(ctx context.Context, eventId uint) (ConditionEventResponse, error) {
	var event ConditionEvent
	var err error
	event, err = getConditionEvent(service, eventId)
	return event.toResponse(), err
}

// encore:api auth method=PATCH path=/condition-events/:eventId
func (service *Service) PatchDate(ctx context.Context, eventId uint, params ConditionEventRequestParams) error {
	return &errs.Error{Code: errs.NotFound, Message: "Endpoint doesn't exist."}
}

// encore:api auth method=DELETE path=/condition-events/:eventId
func (service *Service) DeleteConditionEvent(ctx context.Context, eventId uint) error {
	var event = &ConditionEvent{ID: eventId}
	return event.delete(service)
}
