package pollen

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestData(t *testing.T) DWD {
	var content []byte
	var err error
	content, err = os.ReadFile("./mockResponse.json")
	assert.NoError(t, err)
	var dwd DWD
	dwd, err = parseAPIToDWD(content)
	assert.NoError(t, err)
	return dwd
}

func TestParseRespToDWD(t *testing.T) {
	var dwd DWD = getTestData(t)
	assert.Len(t, dwd.Content, 27)
	assert.Equal(t, "2024-05-31 11:00 Uhr", dwd.LastUpdate)
	assert.Equal(t, "2024-06-01 11:00 Uhr", dwd.NextUpdate)
	var region DWDPollenRegion
	region = dwd.Content[0]
	assert.Equal(t, 10, region.RegionID)
	assert.Equal(t, 11, region.PartregionID)
	var pollen DWDPollen
	pollen = region.Pollen
	assert.Equal(t, "0", pollen.Beifuss.Today)
	assert.Equal(t, "0-1", pollen.Roggen.Today)
	assert.Equal(t, "0-1", pollen.Graeser.Today)
	assert.Equal(t, "0", pollen.Hasel.Today)
	assert.Equal(t, "0", pollen.Esche.Today)
	assert.Equal(t, "0", pollen.Erle.Today)
	assert.Equal(t, "0", pollen.Birke.Today)

}

func TestAPICall(t *testing.T) {
	t.Skip()
	_, err := callDWDAPI()
	assert.NoError(t, err)
}

func TestGetKarlsruheData(t *testing.T) {
	var dwd DWD = getTestData(t)

	var err error
	var pollen DWDPollen
	pollen, err = getKarlsruheData(dwd, Oberrhein)
	assert.NoError(t, err)
	assert.Equal(t, "1", pollen.Graeser.Today)
}

func TestDWDStringToDate(t *testing.T) {
	var timestring string = "2024-06-01 11:00 Uhr"

	var time time.Time
	var err error
	time, err = dwdStringToDate(timestring)
	assert.NoError(t, err)
	assert.Equal(t, "June", time.Month().String())
}
