package hista

import (
	"context"
	"slices"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/pollen"
	"encore.dev/cron"
	"encore.dev/types/uuid"
)

type PollenEventListResponse struct {
	Pollens []PollenEventResponse `json:"pollens"`
}

type PollenEventResponse struct {
	Date    time.Time        `json:"date"`
	Pollens []PollenResponse `json:"pollens"`
}

type PollenResponse struct {
	Type            string `json:"type"`
	Intensity       uint   `json:"intensity"`
	IntensityString string `json:"intensityString"`
}

func toPollenEventListResponse(events []pollen.PollenEvent) PollenEventListResponse {
	var resp = []PollenEventResponse{}
	for _, e := range events {
		var pollens = []PollenResponse{}
		for _, p := range e.Pollens {
			pollens = append(pollens, PollenResponse{
				Type:      string(p.Type),
				Intensity: uint(p.Intensity),
			})
		}
		resp = append(resp, PollenEventResponse{
			Date:    e.CreatedAt,
			Pollens: pollens,
		})
	}
	slices.SortFunc(resp, func(a, b PollenEventResponse) int {
		return b.Date.Compare(a.Date)
	})
	return PollenEventListResponse{Pollens: resp}
}

// encore:api auth method=GET path=/piid/:piid/pollen
func (service *Service) ListPollens(ctx context.Context, piid uuid.UUID) (PollenEventListResponse, error) {
	events, err := service.pollens.ListPollen(ctx)
	if err != nil {
		return PollenEventListResponse{}, errors.MapError(err)
	}
	return toPollenEventListResponse(events), nil
}

var _ = cron.NewJob("pollen", cron.JobConfig{
	Title:    "Get pollen for today",
	Every:    12 * cron.Hour,
	Endpoint: UpdatePollen,
})

// encore:api private method=POST path=/pollen
func (service *Service) UpdatePollen(ctx context.Context) error {
	return errors.MapError(service.pollens.CreatePollen(ctx))
}
