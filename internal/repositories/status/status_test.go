package status

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTest(t *testing.T) *StatusRepo {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.Exec("PRAGMA foreign_keys = ON")
	err = db.AutoMigrate(
		&entity.Status{},
		&entity.MorningStatus{},
		&entity.EveningStatus{},
	)
	assert.NoError(t, err)
	return NewStatusRepo(db)
}

func TestFirst(t *testing.T) {
	repo := initTest(t)

	morningFitness := 3
	repo.db.Create(&entity.Status{Date: time.Now(), MorningFitness: &morningFitness})
	var status = &entity.Status{ID: 1}
	err := repo.First(status)
	assert.NoError(t, err)
	assert.Equal(t, &morningFitness, status.MorningFitness)
}

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		morningFitness int
		eveningFitness int
	}{
		"only date": {},
		"morning":   {morningFitness: 8},
		"evening":   {morningFitness: 8},
		"both":      {morningFitness: 8, eveningFitness: 3},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo := initTest(t)
			var status = &entity.Status{Date: time.Now(), MorningFitness: &test.morningFitness, EveningFitness: &test.eveningFitness}

			err := repo.Create(status)
			assert.NoError(t, err)
			var result entity.Status
			repo.db.First(&result)

			assert.Equal(t, &test.eveningFitness, result.EveningFitness)
			assert.Equal(t, &test.morningFitness, result.MorningFitness)
		})
	}

}

func TestUpdateStatus(t *testing.T) {
	date := time.Date(2023, 12, 11, 10, 9, 8, 6, time.UTC)
	initDate := time.Date(2022, 12, 11, 10, 9, 8, 6, time.UTC)

	changeFitness := 3
	tests := map[string]struct {
		date           time.Time
		morningFitness *int
		eveningFitness *int
	}{
		"only date": {date: date},
		"morning":   {morningFitness: &changeFitness},
		"evening":   {eveningFitness: &changeFitness},
		"both":      {morningFitness: &changeFitness, eveningFitness: &changeFitness},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo := initTest(t)
			morningFitness := 1
			eveningFitness := 5
			var status = &entity.Status{Date: initDate, MorningFitness: &morningFitness, EveningFitness: &eveningFitness}
			repo.Create(status)

			status.MorningFitness = test.morningFitness
			status.EveningFitness = test.eveningFitness
			status.Date = test.date
			err := repo.Update(status)
			assert.NoError(t, err)

			var result entity.Status
			repo.db.First(&result)
			if test.date.IsZero() {
				assert.True(t, result.Date.Equal(initDate))
			} else {
				assert.True(t, result.Date.Equal(date))
			}
			if test.eveningFitness == nil {
				assert.Equal(t, &eveningFitness, result.EveningFitness)
			} else {
				assert.Equal(t, test.eveningFitness, result.EveningFitness)
			}
			if test.morningFitness == nil {
				assert.Equal(t, &morningFitness, result.MorningFitness)
			} else {
				assert.Equal(t, test.morningFitness, result.MorningFitness)
			}
		})
	}
}

func TestFind(t *testing.T) {
	repo := initTest(t)

	var statuses []entity.Status
	statuses = repo.Find()
	assert.Len(t, statuses, 0)

	status := entity.Status{Date: time.Now()}
	repo.db.Create(&status)
	statuses = repo.Find()
	assert.Len(t, statuses, 1)
}

func TestDelete(t *testing.T) {
	repo := initTest(t)

	morningFitness := 1
	status := &entity.Status{Date: time.Now(), MorningFitness: &morningFitness}
	err := repo.db.Create(status).Error
	assert.NoError(t, err)
	err = repo.Delete(status)
	assert.NoError(t, err)

	var rows int64
	repo.db.Find(&entity.Statuses{}).Count(&rows)
	assert.EqualValues(t, 0, rows)
	repo.db.Find(&entity.MorningStatus{}).Count(&rows)
	assert.EqualValues(t, 0, rows)
}

func TestDeleteNotFound(t *testing.T) {
	repo := initTest(t)
	status := &entity.Status{ID: 1}
	err := repo.Delete(status)
	assert.EqualError(t, err, "not_found: not found")
}

func TestFindForDate(t *testing.T) {
	repo := initTest(t)

	status := &entity.Status{ID: 1, Date: time.Now()}
	repo.db.Create(status)

	found, exists := repo.FindForDate(time.Now())
	assert.True(t, exists)
	assert.EqualValues(t, 1, found.ID)
}
