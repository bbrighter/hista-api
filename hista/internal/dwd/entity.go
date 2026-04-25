package dwd

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
