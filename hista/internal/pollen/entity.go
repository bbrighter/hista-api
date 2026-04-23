package pollen

import (
	"time"
)

type PollenEvent struct {
	ID        uint
	CreatedAt time.Time
	Pollens   []Pollen
}

type Pollen struct {
	ID            uint
	PollenEventID uint
	Type          PollenType
	Intensity     PollenIntensity
}

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
