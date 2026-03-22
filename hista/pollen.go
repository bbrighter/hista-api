package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/cron"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/pollen
func (service *Service) ListPollens(ctx context.Context, piid uuid.UUID) (entity.PollenEventsResponse, error) {
	events, err := service.pollens.List(ctx)
	return events.ToResponse(), errors.MapError(err)
}

var _ = cron.NewJob("pollen", cron.JobConfig{
	Title:    "Get pollen for today",
	Every:    12 * cron.Hour,
	Endpoint: UpdatePollen,
})

// encore:api private method=POST path=/pollen
func (service *Service) UpdatePollen(ctx context.Context) error {
	return errors.MapError(service.pollens.Create(ctx))
}
