package status

import (
	"context"
	"testing"
	"time"

	"encore.app/entity"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var GUID = uuid.FromStringOrNil("cf0d4408-8db5-4572-b5d9-4ed873d1341f")

func initTest(t *testing.T) (*StatusRepo, context.Context) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.Exec("PRAGMA foreign_keys = ON")
	err = db.AutoMigrate(
		&entity.Status{},
	)
	assert.NoError(t, err)
	ctx := context.WithValue(t.Context(), "piid", GUID)
	return NewStatusRepo(db), ctx
}

func TestFirst(t *testing.T) {
	repo, ctx := initTest(t)

	morningFitness := 3
	repo.db.Create(&entity.Status{Date: time.Now(), MorningFitness: &morningFitness, PIID: GUID})
	status, err := repo.First(ctx, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, &morningFitness, status.MorningFitness)
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
			repo, ctx := initTest(t)
			var status = &entity.Status{Date: time.Now(), MorningFitness: &test.morningFitness, EveningFitness: &test.eveningFitness, PIID: GUID}

			err := repo.Create(ctx, status)
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
			repo, ctx := initTest(t)
			morningFitness := 1
			eveningFitness := 5
			var status = &entity.Status{Date: initDate, MorningFitness: &morningFitness, EveningFitness: &eveningFitness, PIID: GUID}
			repo.Create(ctx, status)

			status.MorningFitness = test.morningFitness
			status.EveningFitness = test.eveningFitness
			status.Date = test.date
			err := repo.Update(ctx, status)
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
	repo, ctx := initTest(t)

	statuses, err := repo.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, statuses, 0)

	status := entity.Status{Date: time.Now(), PIID: GUID}
	repo.db.Create(&status)
	statuses, err = repo.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, statuses, 1)
}

func TestDelete(t *testing.T) {
	repo, ctx := initTest(t)

	morningFitness := 1
	status := &entity.Status{Date: time.Now(), MorningFitness: &morningFitness, PIID: GUID}
	err := repo.db.Create(status).Error
	assert.NoError(t, err)
	err = repo.Delete(ctx, status.ID)
	assert.NoError(t, err)

	var rows int64
	repo.db.Find(&entity.Statuses{}).Count(&rows)
	assert.EqualValues(t, 0, rows)
	repo.db.Find(&entity.MorningStatus{}).Count(&rows)
	assert.EqualValues(t, 0, rows)
}

func TestDeleteNotFound(t *testing.T) {
	repo, ctx := initTest(t)
	err := repo.Delete(ctx, 1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFindForDate(t *testing.T) {
	repo, ctx := initTest(t)

	status := &entity.Status{ID: 1, Date: time.Now(), PIID: GUID}
	repo.db.Create(status)

	found, exists := repo.FindForDate(ctx, time.Now())
	assert.True(t, exists)
	assert.EqualValues(t, 1, found.ID)
}
