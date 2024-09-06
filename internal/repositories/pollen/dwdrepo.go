package pollen

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

type DWDRepo struct {
	query    func() (entity.DWD, error)
	response *entity.DWD
}

func NewDWDRepo() *DWDRepo {
	return &DWDRepo{
		query:    queryDWDAPI,
		response: new(entity.DWD),
	}
}

func (dwd *DWDRepo) UseTestQuery(t *testing.T) {
	getTestData := func() (dwd entity.DWD, err error) {
		var content []byte
		content, err = os.ReadFile("./mockResponse.json")
		assert.NoError(t, err)
		err = json.Unmarshal(content, &dwd)
		assert.NoError(t, err)
		return dwd, nil
	}
	dwd.query = getTestData
}

func queryDWDAPI() (dwd entity.DWD, err error) {
	const url string = "https://opendata.dwd.de/climate_environment/health/alerts/s31fg.json"
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
