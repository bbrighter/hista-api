package api

import (
	"context"

	entity "encore.app/entity"
	"encore.dev/cron"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/pollen
func (service *Service) ListPollens(ctx context.Context, piid uuid.UUID) (entity.PollenEventsResponse, error) {
	events := service.pollens.List()
	return events.ToResponse(), nil
}

var _ = cron.NewJob("pollen", cron.JobConfig{
	Title:    "Get pollen for today",
	Every:    12 * cron.Hour,
	Endpoint: UpdatePollen,
})

// encore:api private
func (service *Service) UpdatePollen(ctx context.Context) error {
	return service.pollens.Create()
}
