package pollen

import (
	"context"
	"time"

	"encore.dev/cron"
)

var _ = cron.NewJob("pollen", cron.JobConfig{
	Title:    "Get pollen for today",
	Every:    12 * cron.Hour,
	Endpoint: UpdatePollen,
})

// encore:api private
func (service *Service) UpdatePollen(ctx context.Context) error {
	var err error
	var resp []byte
	resp, err = callDWDAPI()
	if err != nil {
		return err
	}
	var dwd DWD
	dwd, err = parseAPIToDWD(resp)
	if err != nil {
		return err
	}
	var lastUpdatedAt time.Time
	lastUpdatedAt, err = dwdStringToDate(dwd.LastUpdate)
	if err != nil {
		return err
	}
	var pollen Pollen
	pollen, err = dwd.toPollen()
	if err != nil {
		return err
	}
	return pollen.writeToDatabase(service, lastUpdatedAt)
}
