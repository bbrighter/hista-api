package dwdPollen

import (
	"encore.app/hista/internal/dwd"
	"encore.app/hista/internal/pollen"
)

func PollenIntensity(intensity dwd.DWDPollenIntensity) pollen.PollenIntensity {
	switch intensity.Today {
	case "0":
		return pollen.NoPollen
	case "0-1":
		return pollen.NoToSmallPollen
	case "1":
		return pollen.SmallPollen
	case "1-2":
		return pollen.SmallToMediumPollen
	case "2":
		return pollen.MediumPollen
	case "2-3":
		return pollen.MediumToHighPollen
	default:
		return pollen.HighPollen
	}
}

func dwdToPollens(dwd dwd.DWDPollen) []pollen.Pollen {
	return []pollen.Pollen{
		{Type: pollen.Ambrosia, Intensity: PollenIntensity(dwd.Ambrosia)},
		{Type: pollen.Roggen, Intensity: PollenIntensity(dwd.Roggen)},
		{Type: pollen.Erle, Intensity: PollenIntensity(dwd.Erle)},
		{Type: pollen.Beifuss, Intensity: PollenIntensity(dwd.Beifuss)},
		{Type: pollen.Birke, Intensity: PollenIntensity(dwd.Birke)},
		{Type: pollen.Graeser, Intensity: PollenIntensity(dwd.Graeser)},
		{Type: pollen.Hasel, Intensity: PollenIntensity(dwd.Hasel)},
		{Type: pollen.Esche, Intensity: PollenIntensity(dwd.Esche)},
	}
}
