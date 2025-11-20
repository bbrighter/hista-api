package hista

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/headaches
func (service *Service) ListHeadaches(ctx context.Context, piid uuid.UUID) (entity.HeadachesResponse, error) {
	headaches, err := service.headaches.List(ctx)
	return headaches.ToResp(), errors.MapError(err)
}

// encore:api auth method=GET path=/piid/:piid/headaches/:id
func (service *Service) GetHeadache(ctx context.Context, piid uuid.UUID, id uint) (entity.HeadacheResponse, error) {
	headache, err := service.headaches.Get(ctx, id)
	return headache.ToResp(), errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/headaches/:id
func (service *Service) DeleteHeadache(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.headaches.Delete(ctx, id))
}

type PostHeadacheParams struct {
	Date     time.Time               `json:"date"`
	Severity entity.HeadacheSeverity `json:"severity"`
}

// encore:api auth method=POST path=/piid/:piid/headaches
func (service *Service) PostHeadache(ctx context.Context, piid uuid.UUID, params PostHeadacheParams) (entity.IDResponse, error) {
	id, err := service.headaches.Create(ctx, params.Date, params.Severity)
	return entity.IDResponse{ID: id}, errors.MapError(err)
}

type PatchHeadacheDateParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id/date
func (service *Service) PatchHeadacheDate(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadacheDateParams) error {
	return errors.MapError(service.headaches.PatchDate(ctx, id, params.Date))
}

type PatchHeadacheSeverityParams struct {
	Severity entity.HeadacheSeverity `json:"severity"`
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id/severity
func (service *Service) PatchHeadacheSeverity(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadacheSeverityParams) error {
	return errors.MapError(service.headaches.PatchSeverity(ctx, id, params.Severity))
}

type PatchHeadachePositionsParams struct {
	Positions entity.HeadachePositions `json:"positions"`
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id/positions
func (service *Service) PatchHeadachePositions(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadachePositionsParams) error {
	return errors.MapError(service.headaches.PatchPositions(ctx, id, params.Positions))
}

type PatchHeadacheSymptomsParams struct {
	Symptoms entity.HeadacheSymptoms `json:"symptoms"`
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id/symptoms
func (service *Service) PatchHeadacheSymptoms(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadacheSymptomsParams) error {
	return errors.MapError(service.headaches.PatchSymptoms(ctx, id, params.Symptoms))
}

type PatchHeadacheTypesParams struct {
	Types entity.HeadacheTypes `json:"types"`
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id/types
func (service *Service) PatchHeadacheTypes(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadacheTypesParams) error {
	return errors.MapError(service.headaches.PatchTypes(ctx, id, params.Types))
}

type PatchHeadacheDescriptionParams struct {
	Description string `json:"description"`
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id/description
func (service *Service) PatchHeadacheDescription(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadacheDescriptionParams) error {
	return errors.MapError(service.headaches.PatchDescription(ctx, id, params.Description))
}
