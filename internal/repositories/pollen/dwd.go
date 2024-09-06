package pollen

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/errs"
)

type PartRegion int

func (repo *DWDRepo) GetKarlsruheData() (pollen entity.DWDPollen, err error) {
	dwd, err := repo.query()
	repo.response = &dwd
	if err != nil {
		return pollen, err
	}
	const Oberrhein PartRegion = 111
	var index int
	var found bool = false
	for i, con := range dwd.Content {
		if con.PartregionID == int(Oberrhein) {
			index = i
			found = true
		}
	}
	if !found {
		return pollen, errors.ErrorNotFound
	}
	return dwd.Content[index].Pollen, nil
}

func (repo *DWDRepo) DwdStringToDate() (t time.Time, err error) {
	if repo.response == nil {
		return t, errors.NewError("Response is nil", errs.Internal)
	}
	var loc *time.Location
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Now(), err
	}
	return time.ParseInLocation("2006-01-02 15:04 Uhr", repo.response.LastUpdate, loc)
}
