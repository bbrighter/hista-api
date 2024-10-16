package status

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
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

func TestCreate(t *testing.T) {
	repo := initTest(t)

	status := &entity.Status{}
	err := repo.Save(status)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, status.ID)
}

func TestDelete(t *testing.T) {
	repo := initTest(t)

	status := &entity.Status{Date: time.Now(), Morning: &entity.MorningStatus{Fitness: entity.Bad}}
	err := repo.Save(status)
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
	status := &entity.Status{ID: 1, Morning: &entity.MorningStatus{ID: 1, StatusID: 1}}
	err := repo.Delete(status)
	assert.EqualError(t, err, "not_found: not found")
}

func TestFindForDate(t *testing.T) {
	repo := initTest(t)

	status := &entity.Status{ID: 1, Date: time.Now(), Morning: &entity.MorningStatus{Fitness: 4}}
	repo.db.Create(status)

	found, exists := repo.FindForDate(time.Now())
	assert.True(t, exists)
	assert.NotNil(t, found.Morning)
	assert.Nil(t, found.Evening)
}
