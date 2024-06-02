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
	Roggen   PollenIntensitiy `json:"Roggen"`
	Ambrosia PollenIntensitiy `json:"Ambrosia"`
	Erle     PollenIntensitiy `json:"Erle"`
	Beifuss  PollenIntensitiy `json:"Beifuss"`
	Birke    PollenIntensitiy `json:"Birke"`
	Graeser  PollenIntensitiy `json:"Graeser"`
	Hasel    PollenIntensitiy `json:"Hasel"`
	Esche    PollenIntensitiy `json:"Esche"`
}

type PollenIntensitiy struct {
	Today string `json:"today"`
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

func getKarlsruheData(dwd DWD, partregion PartRegion) (DWDPollen, error) {
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
