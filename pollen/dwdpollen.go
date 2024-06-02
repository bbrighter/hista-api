package pollen

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"encore.app/errors"
)

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
		return No
	case "0-1":
		return NoToSmall
	case "1":
		return Small
	case "1-2":
		return SmallToMedium
	case "2":
		return Medium
	case "2-3":
		return MediumToHigh
	default:
		return High
	}
}

func callDWDAPI() ([]byte, error) {
	const url string = "https://opendata.dwd.de/climate_environment/health/alerts/s31fg.json"
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	return io.ReadAll(resp.Body)
}

func parseAPIToDWD(resp []byte) (DWD, error) {
	var pollen DWD
	var err error
	err = json.Unmarshal(resp, &pollen)
	return pollen, err
}

type PartRegion int

const (
	Oberrhein PartRegion = 111
)

func (dwd DWD) getKarlsruheData(partregion PartRegion) (DWDPollen, error) {
	var index int
	var found bool = false
	for i, con := range dwd.Content {
		if con.PartregionID == int(partregion) {
			index = i
			found = true
		}
	}
	if !found {
		return DWDPollen{}, errors.ErrorNotFound
	}
	return dwd.Content[index].Pollen, nil
}

func dwdStringToDate(datestring string) (time.Time, error) {
	var loc *time.Location
	var err error
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Now(), err
	}
	return time.ParseInLocation("2006-01-02 15:04 Uhr", datestring, loc)
}
