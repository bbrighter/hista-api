package dwd

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"encore.app/errors"
)

const DWD_URL = "https://opendata.dwd.de/climate_environment/health/alerts/s31fg.json"

type dwdQuerier interface {
	query() func() (DWD, error)
}

type DWDClient struct {
	query func() (DWD, error)
}

func NewDWDClient() *DWDClient {
	queryDWDAPI := func() (dwd DWD, err error) {
		const url string = DWD_URL
		resp, err := http.Get(url)
		if err != nil {
			return dwd, err
		}
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return dwd, err
		}
		err = json.Unmarshal(bytes, &dwd)
		return dwd, err
	}

	return &DWDClient{
		query: queryDWDAPI,
	}
}

func (client *DWDClient) GetKarlsruheData() (DWDPollen, time.Time, error) {
	dwd, err := client.query()
	if err != nil {
		return DWDPollen{}, time.Time{}, err
	}

	updatedAt, err := dwd.extractDate()
	if err != nil {
		return DWDPollen{}, time.Time{}, err
	}

	const Oberrhein int = 111
	var index int
	var found bool = false
	for i, con := range dwd.Content {
		if con.PartregionID == Oberrhein {
			index = i
			found = true
		}
	}
	if !found {
		return DWDPollen{}, time.Time{}, errors.ErrorNotFound
	}
	return dwd.Content[index].Pollen, updatedAt, nil
}

func (dwd DWD) extractDate() (t time.Time, err error) {
	var loc *time.Location
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Now(), err
	}
	return time.ParseInLocation("2006-01-02 15:04 Uhr", dwd.LastUpdate, loc)
}
