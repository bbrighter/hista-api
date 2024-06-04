package pollen

import (
	"context"
	"time"
)

type PollenEventsResponse struct {
	Pollens []PollenEventResponse `json:"pollens"`
}

type PollenEventResponse struct {
	Date    time.Time        `json:"date"`
	Pollens []PollenResponse `json:"pollens"`
}

type PollenResponse struct {
	Type            PollenType      `json:"type"`
	Intensity       PollenIntensity `json:"intensity"`
	IntensityString string          `json:"intensityString"`
}

// encore:api auth method=GET path=/pollen
func (service *Service) GetPollens(ctx context.Context) (PollenEventsResponse, error) {
	var events PollenEvents = findPollens(service)
	return events.PollenEventsResponse(), nil
}
