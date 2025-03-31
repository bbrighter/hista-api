package api

import (
	"context"
	"time"

	"encore.app/entity"
)

// encore:api auth method=GET path=/headaches
func (service *Service) GetHeadaches(ctx context.Context) (entity.HeadachesResponse, error) {
	headaches := service.headaches.List()
	return headaches.ToResp(), nil
}

// encore:api auth method=GET path=/headaches/:id
func (service *Service) GetHeadache(ctx context.Context, id uint) (entity.HeadacheResponse, error) {
	headache, err := service.headaches.Get(id)
	return headache.ToResp(), err
}

// encore:api auth method=DELETE path=/headaches/:id
func (service *Service) DeleteHeadache(ctx context.Context, id uint) error {
	return service.headaches.Delete(id)
}

type PostHeadacheParams struct {
	Date     time.Time               `json:"date"`
	Severity entity.HeadacheSeverity `json:"severity"`
}

// encore:api auth method=POST path=/headaches
func (service *Service) PostHeadache(ctx context.Context, params PostHeadacheParams) (entity.IDResponse, error) {
	id, err := service.headaches.Create(params.Date, params.Severity)
	return entity.IDResponse{ID: id}, err
}

type PatchHeadacheDateParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=PATCH path=/headaches/:id/date
func (service *Service) PatchHeadacheDate(ctx context.Context, id uint, params PatchHeadacheDateParams) error {
	return service.headaches.PatchDate(id, params.Date)
}

type PatchHeadacheSeverityParams struct {
	Severity entity.HeadacheSeverity `json:"severity"`
}

// encore:api auth method=PATCH path=/headaches/:id/severity
func (service *Service) PatchHeadacheSeverity(ctx context.Context, id uint, params PatchHeadacheSeverityParams) error {
	return service.headaches.PatchSeverity(id, params.Severity)
}

type PatchHeadachePositionsParams struct {
	Positions entity.HeadachePositions `json:"positions"`
}

// encore:api auth method=PATCH path=/headaches/:id/positions
func (service *Service) PatchHeadachePositions(ctx context.Context, id uint, params PatchHeadachePositionsParams) error {
	return service.headaches.PatchPositions(id, params.Positions)
}

type PatchHeadacheSymptomsParams struct {
	Symptoms entity.HeadacheSymptoms `json:"symptoms"`
}

// encore:api auth method=PATCH path=/headaches/:id/symptoms
func (service *Service) PatchHeadacheSymptoms(ctx context.Context, id uint, params PatchHeadacheSymptomsParams) error {
	return service.headaches.PatchSymptoms(id, params.Symptoms)
}

type PatchHeadacheTypesParams struct {
	Types entity.HeadacheTypes `json:"types"`
}

// encore:api auth method=PATCH path=/headaches/:id/types
func (service *Service) PatchHeadacheTypes(ctx context.Context, id uint, params PatchHeadacheTypesParams) error {
	return service.headaches.PatchTypes(id, params.Types)
}

type PatchHeadacheDescriptionParams struct {
	Description string `json:"description"`
}

// encore:api auth method=PATCH path=/headaches/:id/description
func (service *Service) PatchHeadacheDescription(ctx context.Context, id uint, params PatchHeadacheDescriptionParams) error {
	return service.headaches.PatchDescription(id, params.Description)
}
