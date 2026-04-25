package hista

import (
	"context"
	"fmt"
	"sort"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/headaches"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type HeadacheResponse struct {
	ID          uint      `json:"id"`
	Date        time.Time `json:"date"`
	Severity    uint8     `json:"severity"`
	Types       []string  `json:"types" encore:"optional"`
	Positions   []string  `json:"positions" encore:"optional"`
	Symptoms    []string  `json:"symptoms" encore:"optional"`
	Description string    `json:"description"`
}

func toStringSlice[T ~string](input []T) []string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = string(v)
	}
	return result
}

func toHeadacheResponse(h *headaches.Headache) HeadacheResponse {
	return HeadacheResponse{
		ID:          h.ID,
		Date:        h.Date,
		Severity:    h.Severity,
		Types:       toStringSlice(h.Types),
		Positions:   toStringSlice(h.Positions),
		Symptoms:    toStringSlice(h.Symptoms),
		Description: h.Description,
	}
}

type HeadacheListResponse struct {
	Headaches []HeadacheResponse `json:"headaches"`
}

func toHeadacheListResponse(hs headaches.Headaches) HeadacheListResponse {
	var headaches = []HeadacheResponse{}
	for _, h := range hs {
		headaches = append(headaches, toHeadacheResponse(h))
	}
	sort.Slice(headaches, func(i, j int) bool {
		return headaches[i].Date.After(headaches[j].Date)
	})
	return HeadacheListResponse{Headaches: headaches}
}

// encore:api auth method=GET path=/piid/:piid/headaches
func (service *Service) ListHeadaches(ctx context.Context, piid uuid.UUID) (HeadacheListResponse, error) {
	headaches, err := service.headaches.ListHeadaches(ctx)
	if err != nil {
		return HeadacheListResponse{}, errors.MapError(err)
	}
	return toHeadacheListResponse(headaches), nil
}

// encore:api auth method=GET path=/piid/:piid/headaches/:id
func (service *Service) GetHeadache(ctx context.Context, piid uuid.UUID, id uint) (HeadacheResponse, error) {
	headache, err := service.headaches.FirstHeadache(ctx, id)
	if err != nil {
		return HeadacheResponse{}, errors.MapError(err)
	}
	return toHeadacheResponse(headache), nil
}

// encore:api auth method=DELETE path=/piid/:piid/headaches/:id
func (service *Service) DeleteHeadache(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.headaches.DeleteHeadache(ctx, id))
}

type PostHeadacheParams struct {
	Date     time.Time `json:"date"`
	Severity uint8     `json:"severity"`
}

// encore:api auth method=POST path=/piid/:piid/headaches
func (service *Service) PostHeadache(ctx context.Context, piid uuid.UUID, params PostHeadacheParams) (IDResponse, error) {
	id, err := service.headaches.CreateHeadache(ctx, params.Date, params.Severity)
	return IDResponse{ID: id}, errors.MapError(err)
}

type PatchHeadacheParams struct {
	Date        option.Option[time.Time]                   `json:"date"`
	Description option.Option[string]                      `json:"description"`
	Severity    option.Option[uint8]                       `json:"severity"`
	Types       option.Option[headaches.HeadacheTypes]     `json:"types"`
	Symptoms    option.Option[headaches.HeadacheSymptoms]  `json:"symptoms"`
	Positions   option.Option[headaches.HeadachePositions] `json:"positions"`
}

func (p PatchHeadacheParams) Validate() error {
	if p.Types.IsSome() {
		for _, typ := range p.Types.MustGet() {
			if !headaches.IsValidHeadacheType(typ) {
				return fmt.Errorf("invalid headache type %s", typ)
			}
		}
	}
	if p.Positions.IsSome() {
		for _, pos := range p.Positions.MustGet() {
			if !headaches.IsValidHeadachePosition(pos) {
				return fmt.Errorf("invalid headache type %s", pos)
			}
		}
	}
	if p.Symptoms.IsSome() {
		for _, sym := range p.Symptoms.MustGet() {
			if !headaches.IsValidSymptom(sym) {
				return fmt.Errorf("invalid headache symptom %s", sym)
			}
		}
	}
	return nil
}

// encore:api auth method=PATCH path=/piid/:piid/headaches/:id
func (s *Service) PatchHeadache(ctx context.Context, piid uuid.UUID, id uint, params PatchHeadacheParams) error {
	err := s.headaches.UpdateHeadache(ctx, id, headaches.HeadacheUpdateParams{
		Date:        params.Date.PtrOrNil(),
		Positions:   params.Positions.PtrOrNil(),
		Severity:    params.Severity.PtrOrNil(),
		Symptoms:    params.Symptoms.PtrOrNil(),
		Types:       params.Types.PtrOrNil(),
		Description: params.Description.PtrOrNil(),
	})
	return errors.MapError(err)
}
