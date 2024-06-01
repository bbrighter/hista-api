package pollen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToPollen(t *testing.T) {
	var dwd DWD = getTestData(t)

	var pollen Pollen
	var err error
	pollen, err = dwd.toPollen()
	assert.NoError(t, err)
	assert.Equal(t, No, pollen.Hasel)
	assert.Equal(t, No, pollen.Esche)
	assert.Equal(t, Small, pollen.Graeser)
	assert.Equal(t, Medium, pollen.Ambrosia)
	assert.Equal(t, SmallToMedium, pollen.Erle)
	assert.Equal(t, NoToSmall, pollen.Roggen)
	assert.Equal(t, High, pollen.Birke)
	assert.Equal(t, MediumToHigh, pollen.Beifuss)
}

func TestWriteToDatabase(t *testing.T) {
	service := initTest(t)

	var pollen = Pollen{Roggen: High}
	var err error
	err = pollen.writeToDatabase(service, time.Now())
	assert.NoError(t, err)

	var pollens []Pollen
	service.db.Find(&pollens)
	assert.Len(t, pollens, 1)

	err = pollen.writeToDatabase(service, time.Now().Add(-time.Hour))
	assert.NoError(t, err)
	service.db.Find(&pollens)
	assert.Len(t, pollens, 1)

	var newPollen = Pollen{Roggen: No}
	err = newPollen.writeToDatabase(service, time.Now().Add(time.Hour))
	assert.NoError(t, err)
	service.db.Find(&pollens)
	assert.Len(t, pollens, 2)
}
