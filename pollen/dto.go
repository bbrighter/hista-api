package pollen

import (
	"slices"
)

func (event PollenEvent) PollenEventResponse() PollenEventResponse {
	var pollens = []PollenResponse{}
	for _, pol := range event.Pollens {
		pollens = append(pollens, PollenResponse{
			Type:            pol.Type,
			Intensity:       pol.Intensity,
			IntensityString: pol.Intensity.String(),
		})
	}
	return PollenEventResponse{
		Date:    event.CreatedAt,
		Pollens: pollens,
	}
}

func (event PollenEvents) PollenEventsResponse() PollenEventsResponse {
	var resp = []PollenEventResponse{}
	for _, e := range event {
		resp = append(resp, e.PollenEventResponse())
	}
	slices.SortFunc(resp, func(a, b PollenEventResponse) int {
		return a.Date.Compare(b.Date)
	})
	return PollenEventsResponse{Pollens: resp}
}
