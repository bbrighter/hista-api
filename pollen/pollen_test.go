package pollen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToPollen(t *testing.T) {
	var dwd DWD = getTestData(t)

	var pollens Pollens
	var err error
	pollens, err = dwd.Pollens()
	assert.NoError(t, err)
	assert.Len(t, pollens, 8)
	for _, pollen := range pollens {
		switch pollen.Type {
		case Hasel:
			assert.Equal(t, No, pollen.Intensity)
		case Esche:
			assert.Equal(t, No, pollen.Intensity)
		case Graeser:
			assert.Equal(t, Small, pollen.Intensity)
		case Ambrosia:
			assert.Equal(t, Medium, pollen.Intensity)
		case Erle:
			assert.Equal(t, SmallToMedium, pollen.Intensity)
		case Roggen:
			assert.Equal(t, NoToSmall, pollen.Intensity)
		case Birke:
			assert.Equal(t, High, pollen.Intensity)
		case Beifuss:
			assert.Equal(t, MediumToHigh, pollen.Intensity)
		}
	}
}

func TestWriteToDatabase(t *testing.T) {
	service := initTest(t)

	var event = &PollenEvent{Pollens: Pollens{Pollen{Type: Roggen, Intensity: High}}}
	var err error
	err = event.create(service, time.Now())
	assert.NoError(t, err)
	defer service.db.Delete(event)

	var pollens []Pollen
	service.db.Find(&pollens)
	assert.Len(t, pollens, 1)

	err = event.create(service, time.Now().Add(-time.Hour))
	assert.NoError(t, err)
	service.db.Find(&pollens)
	assert.Len(t, pollens, 1)

	// var newPollen = Pollen{Roggen: No}
	var newEvent = &PollenEvent{Pollens: Pollens{Pollen{Type: Roggen, Intensity: No}}}
	err = newEvent.create(service, time.Now().Add(time.Hour))
	defer service.db.Delete(newEvent)
	assert.NoError(t, err)
	service.db.Find(&pollens)
	assert.Len(t, pollens, 2)

}

func TestFindPollenWithSeverity(t *testing.T) {
	service := initTest(t)
	dwd := getTestData(t)
	pollen, _ := dwd.Pollens()
	var pollenEvent = &PollenEvent{Pollens: pollen}
	pollenEvent.create(service, time.Now())
	defer service.db.Delete(&pollenEvent)

	pollens := FindPollenWithSeverity(service.db)
	assert.Len(t, pollens[0].Pollens, 6)
}

func TestFindPollens(t *testing.T) {
	service := initTest(t)
	dwd := getTestData(t)
	pollen, _ := dwd.Pollens()
	var pollenEvent = &PollenEvent{Pollens: pollen}
	pollenEvent.create(service, time.Now())
	defer service.db.Delete(&pollenEvent)

	var events PollenEvents = findPollens(service)
	assert.Len(t, events, 1)
	assert.Len(t, events[0].Pollens, 8)
}
