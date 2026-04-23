package hista

import (
	"context"
	"sort"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/symptoms"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type ConditionEventListResponse struct {
	ConditionEvents []ConditionEventMetaResponse `json:"conditionEvents"`
}
type ConditionEventMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type ConditionEventResponse struct {
	ConditionEventMetaResponse
	Conditions []ConditionResponse `json:"conditions"`
}

type ConditionResponse struct {
	ID       uint            `json:"id"`
	Symptom  SymptomResponse `json:"symptom"`
	Severity uint            `json:"severity"`
}

func toConditionResponse(c symptoms.Condition) ConditionResponse {
	return ConditionResponse{
		ID:       c.ID,
		Symptom:  toSymptomResponse(c.Symptom),
		Severity: uint(c.Severity),
	}
}

func toConditionEventListResponse(events symptoms.ConditionEvents) ConditionEventListResponse {
	var resp = []ConditionEventMetaResponse{}
	for _, s := range events {
		resp = append(resp,
			ConditionEventMetaResponse{
				ID:   s.ID,
				Date: s.Date,
			})
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Date.Sub(resp[j].Date) > 0
	})
	return ConditionEventListResponse{ConditionEvents: resp}
}

func toConditionEventResponse(c symptoms.ConditionEvent) ConditionEventResponse {
	var conditionResp = []ConditionResponse{}
	for _, con := range c.Conditions {
		conditionResp = append(conditionResp, toConditionResponse(con))
	}
	var resp = ConditionEventResponse{
		ConditionEventMetaResponse: ConditionEventMetaResponse{ID: c.ID, Date: c.Date},
		Conditions:                 conditionResp,
	}
	return resp
}

// encore:api auth method=POST path=/piid/:piid/condition-events
func (service *Service) CreateConditionEvent(ctx context.Context, piid uuid.UUID) (ConditionEventResponse, error) {
	event, err := service.syms.CreateConditionEvent(ctx)
	return toConditionEventResponse(event), errors.MapError(err)
}

// encore:api auth method=GET path=/piid/:piid/condition-events
func (service *Service) ListConditionEvents(ctx context.Context, piid uuid.UUID) (ConditionEventListResponse, error) {
	events, err := service.syms.ListConditionEvents(ctx)
	return toConditionEventListResponse(events), errors.MapError(err)
}

// encore:api auth method=GET path=/piid/:piid/condition-events/:eventId
func (service *Service) GetConditionEvent(ctx context.Context, piid uuid.UUID, eventId uint) (ConditionEventResponse, error) {
	event, err := service.syms.GetConditionEvent(ctx, eventId)
	return toConditionEventResponse(event), errors.MapError(err)
}

type ConditionEventRequestParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=PATCH path=/piid/:piid/condition-events/:eventId
func (service *Service) PatchDate(ctx context.Context, piid uuid.UUID, eventId uint, params ConditionEventRequestParams) error {
	if params.Date.IsZero() {
		return errors.ErrorAttributeMustBeSet("date")
	}
	err := service.syms.PatchConditionEvent(ctx, eventId, params.Date)
	return errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/condition-events/:eventId
func (service *Service) DeleteConditionEvent(ctx context.Context, piid uuid.UUID, eventId uint) (SymptomCategoryListResponse, error) {
	cats, err := service.syms.DeleteConditionEvent(ctx, eventId)
	return toSymptomCategoryListResponse(cats), errors.MapError(err)
}

type ConditionRequestParams struct {
	SymptomName option.Option[string] `json:"symptomName" encore:"optional"`
	SymptomID   option.Option[uint]   `json:"symptomId" encore:"optional"`
	CategoryID  option.Option[uint]   `json:"categoryId" encore:"optional"`
}

type PostConditionResponse struct {
	Condition ConditionResponse                          `json:"condition"`
	Symptoms  option.Option[SymptomCategoryListResponse] `json:"symptoms" encore:"optional"`
}

// encore:api auth method=POST path=/piid/:piid/condition-events/:eventId/conditions
func (service *Service) PostCondition(ctx context.Context, piid uuid.UUID, eventId uint, params ConditionRequestParams) (PostConditionResponse, error) {
	var err error = errors.BadRequest("params empty")
	var id uint
	var cats symptoms.SymptomCategories
	if symptomId, ok := params.SymptomID.Get(); ok {
		id, err = service.syms.CreateConditionById(ctx, eventId, symptomId)
	} else if params.SymptomName.IsSome() && params.CategoryID.IsSome() {
		symptomName := params.SymptomName.GetOrElse("")
		categoryId := params.CategoryID.GetOrElse(0)
		id, cats, err = service.syms.CreateConditionWithNewIngredient(ctx, eventId, symptomName, categoryId)
	}
	if err != nil {
		return PostConditionResponse{}, errors.MapError(err)
	}

	var symptomOpts option.Option[SymptomCategoryListResponse]
	if len(cats) > 0 {
		symptomOpts = option.Some(toSymptomCategoryListResponse(cats))
	} else {
		symptomOpts = option.None[SymptomCategoryListResponse]()
	}

	return PostConditionResponse{
		Condition: ConditionResponse{ID: id},
		Symptoms:  symptomOpts,
	}, nil

}
