package pollen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func initTest(t *testing.T) *DWDRepo {
	repo := NewDWDRepo()
	repo.UseTestQuery(t)
	return repo
}

func TestGetKarlsruheData(t *testing.T) {
	repo := initTest(t)
	dwd, err := repo.GetKarlsruheData()

	assert.NoError(t, err)
	assert.Equal(t, "2-3", dwd.Beifuss.Today)
	assert.Equal(t, "0-1", dwd.Roggen.Today)
	assert.Equal(t, "1", dwd.Graeser.Today)
	assert.Equal(t, "0", dwd.Hasel.Today)
	assert.Equal(t, "0", dwd.Esche.Today)
	assert.Equal(t, "1-2", dwd.Erle.Today)
	assert.Equal(t, "3", dwd.Birke.Today)
	assert.Equal(t, "2", dwd.Ambrosia.Today)
}

// func TestAPICall(t *testing.T) {
// 	_, err := callDWDAPI()
// 	assert.NoError(t, err)
// }

// func TestGetKarlsruheData(t *testing.T) {
// 	var dwd DWD
// 	dwd, _ = getTestData()

// 	var err error
// 	var pollen Pollen
// 	pollen, err = dwd.getKarlsruheData(Oberrhein)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "1", pollen.Graeser.Today)
// }

func TestDWDStringToDate(t *testing.T) {
	repo := initTest(t)
	dwd, err := repo.query()
	assert.NoError(t, err)
	repo.response = &dwd

	time, err := repo.DwdStringToDate()
	assert.NoError(t, err)
	assert.Equal(t, "May", time.Month().String())
}
