package entity

import (
	"slices"
	"time"
)

type PollenIntensity uint

const (
	NoPollen            PollenIntensity = 1
	NoToSmallPollen     PollenIntensity = 2
	SmallPollen         PollenIntensity = 3
	SmallToMediumPollen PollenIntensity = 4
	MediumPollen        PollenIntensity = 5
	MediumToHighPollen  PollenIntensity = 6
	HighPollen          PollenIntensity = 7
)

type PollenType string

const (
	Roggen   PollenType = "Roggen"
	Ambrosia PollenType = "Ambrosia"
	Erle     PollenType = "Erle"
	Beifuss  PollenType = "Beifuss"
	Birke    PollenType = "Birke"
	Graeser  PollenType = "Gräser"
	Hasel    PollenType = "Hasel"
	Esche    PollenType = "Esche"
)

func (intensity PollenIntensity) String() string {
	var intensityString string
	switch intensity {
	case NoPollen:
		intensityString = "Keine"
	case NoToSmallPollen:
		intensityString = "Keine bis geringe"
	case SmallPollen:
		intensityString = "Geringe"
	case SmallToMediumPollen:
		intensityString = "Geringe bis mittlere"
	case MediumPollen:
		intensityString = "Mittlere"
	case MediumToHighPollen:
		intensityString = "Mittlere bis hohe"
	case HighPollen:
		intensityString = "Hohe"
	}
	return intensityString
}

type PollenEvent struct {
	ID        uint
	CreatedAt time.Time
	Pollens   Pollens
}

type PollenEvents []PollenEvent

type Pollen struct {
	ID            uint
	PollenEventID uint
	Type          PollenType
	Intensity     PollenIntensity
}

type Pollens []Pollen

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

func (event PollenEvent) ToResponse() PollenEventResponse {
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

func (event PollenEvents) ToResponse() PollenEventsResponse {
	var resp = []PollenEventResponse{}
	for _, e := range event {
		resp = append(resp, e.ToResponse())
	}
	slices.SortFunc(resp, func(a, b PollenEventResponse) int {
		return b.Date.Compare(a.Date)
	})
	return PollenEventsResponse{Pollens: resp}
}
