package pollen

import (
	"testing"
	"time"

	"encore.app/hista/entity"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initPollenTest(t *testing.T) *PollenRepo {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&entity.Pollen{}, entity.PollenEvent{})
	assert.NoError(t, err)
	return &PollenRepo{
		db: db,
	}
}

func TestCreate(t *testing.T) {
	repo := initPollenTest(t)

	dwdLastUpdated := time.Date(2000, 12, 12, 0, 0, 0, 0, time.Local)
	var err error
	var inputPollens = entity.Pollens{
		entity.Pollen{
			Type:      entity.Ambrosia,
			Intensity: entity.HighPollen,
		},
		entity.Pollen{
			Type:      entity.Beifuss,
			Intensity: entity.MediumToHighPollen,
		},
	}
	err = repo.Create(inputPollens, dwdLastUpdated)
	assert.NoError(t, err)

	var pollens entity.Pollens
	var events entity.PollenEvents
	repo.db.Find(&pollens)
	assert.Len(t, pollens, 2)
	repo.db.Find(&events)
	assert.Len(t, events, 1)

	var inputPollens2 = entity.Pollens{
		entity.Pollen{
			Type:      entity.Ambrosia,
			Intensity: entity.HighPollen,
		},
		entity.Pollen{
			Type:      entity.Beifuss,
			Intensity: entity.MediumToHighPollen,
		},
	}
	err = repo.Create(inputPollens2, dwdLastUpdated)
	assert.NoError(t, err)
	repo.db.Find(&pollens)
	assert.Len(t, pollens, 2)
	repo.db.Find(&events)
	assert.Len(t, events, 1)

	var inputPollens3 = entity.Pollens{
		entity.Pollen{
			Type:      entity.Ambrosia,
			Intensity: entity.HighPollen,
		},
		entity.Pollen{
			Type:      entity.Beifuss,
			Intensity: entity.MediumToHighPollen,
		},
	}
	err = repo.Create(inputPollens3, time.Now().Add(time.Hour))
	assert.NoError(t, err)
	repo.db.Find(&pollens)
	assert.Len(t, pollens, 4)
	repo.db.Find(&events)
	assert.Len(t, events, 2)
}

func TestFindPollenWithSeverity(t *testing.T) {
	repo := initPollenTest(t)

	var events entity.PollenEvents
	events = repo.FindPollenWithSeverity(1)
	assert.Len(t, events, 0)

	// Fill db
	var input = entity.Pollens{
		entity.Pollen{Intensity: entity.HighPollen},
		entity.Pollen{Intensity: entity.NoPollen},
	}
	err := repo.Create(input, time.Now())
	assert.NoError(t, err)

	events = repo.FindPollenWithSeverity(1)
	assert.Len(t, events, 1)
	pollens := events[0].Pollens
	assert.Len(t, pollens, 1)

	events = repo.FindPollenWithSeverity(0)
	assert.Len(t, events, 1)
	pollens = events[0].Pollens
	assert.Len(t, pollens, 2)
}
