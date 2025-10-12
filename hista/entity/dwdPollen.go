package entity

type DWD struct {
	NextUpdate string            `json:"next_update"`
	LastUpdate string            `json:"last_update"`
	Content    []DWDPollenRegion `json:"content"`
}

type DWDPollenRegion struct {
	RegionID     int       `json:"region_id"`
	PartregionID int       `json:"partregion_id"`
	Pollen       DWDPollen `json:"Pollen"`
}

type DWDPollen struct {
	Roggen   DWDPollenIntensity `json:"Roggen"`
	Ambrosia DWDPollenIntensity `json:"Ambrosia"`
	Erle     DWDPollenIntensity `json:"Erle"`
	Beifuss  DWDPollenIntensity `json:"Beifuss"`
	Birke    DWDPollenIntensity `json:"Birke"`
	Graeser  DWDPollenIntensity `json:"Graeser"`
	Hasel    DWDPollenIntensity `json:"Hasel"`
	Esche    DWDPollenIntensity `json:"Esche"`
}

type DWDPollenIntensity struct {
	Today string `json:"today"`
}

func (intensity DWDPollenIntensity) PollenIntensity() PollenIntensity {
	switch intensity.Today {
	case "0":
		return NoPollen
	case "0-1":
		return NoToSmallPollen
	case "1":
		return SmallPollen
	case "1-2":
		return SmallToMediumPollen
	case "2":
		return MediumPollen
	case "2-3":
		return MediumToHighPollen
	default:
		return HighPollen
	}
}

func (dwd DWDPollen) ToPollen() Pollens {
	return Pollens{
		Pollen{Type: Ambrosia, Intensity: dwd.Ambrosia.PollenIntensity()},
		Pollen{Type: Roggen, Intensity: dwd.Roggen.PollenIntensity()},
		Pollen{Type: Erle, Intensity: dwd.Erle.PollenIntensity()},
		Pollen{Type: Beifuss, Intensity: dwd.Beifuss.PollenIntensity()},
		Pollen{Type: Birke, Intensity: dwd.Birke.PollenIntensity()},
		Pollen{Type: Graeser, Intensity: dwd.Graeser.PollenIntensity()},
		Pollen{Type: Hasel, Intensity: dwd.Hasel.PollenIntensity()},
		Pollen{Type: Esche, Intensity: dwd.Esche.PollenIntensity()},
	}
}
