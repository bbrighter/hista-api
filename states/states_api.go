package states

import (
	"context"
	"time"

	"encore.dev/beta/errs"
)

type StateRequestParams struct {
	Date time.Time `json:"date"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

// encore:api auth method=POST path=/state
func (service *Service) CreateState(ctx context.Context, params StateRequestParams) (IDResponse, error) {
	var state *State = newState(params.Date)
	var err error = state.create(service)
	return IDResponse{ID: state.ID}, err
}

type StatesResponse struct {
	States []StateMetaResponse `json:"states"`
}

type StateMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

// encore:api auth method=GET path=/state
func (service *Service) GetStates(ctx context.Context) (StatesResponse, error) {
	return getStates(service).toResponse(), nil
}

type StateResponse struct {
	ID         uint                `json:"id"`
	Date       time.Time           `json:"date"`
	Conditions []ConditionResponse `json:"conditions"`
}

type ConditionResponse struct {
	ID       uint              `json:"id"`
	Symptom  SymptomResponse   `json:"symptom"`
	Severity ConditionSeverity `json:"severity"`
}

// encore:api auth method=GET path=/state/:stateId
func (service *Service) GetState(ctx context.Context, stateId uint) (StateResponse, error) {
	var state State
	var err error
	state, err = getState(service, stateId)
	return state.toResponse(), err
}

// encore:api auth method=PATCH path=/state/:stateId
func (service *Service) PatchDate(ctx context.Context, stateId uint, params StateRequestParams) error {
	return &errs.Error{Code: errs.NotFound, Message: "Endpoint doesn't exist."}
}

// encore:api auth method=DELETE path=/state/:stateId
func (service *Service) DeleteState(ctx context.Context, stateId uint) error {
	var state = &State{ID: stateId}
	return state.delete(service)
}
