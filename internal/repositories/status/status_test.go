package status

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func TestFirst(t *testing.T) {
	repo := initTest(t)

	repo.db.Create(&entity.Status{Date: time.Now(), Morning: &entity.MorningStatus{Fitness: entity.Bad}})
	var status = &entity.Status{ID: 1}
	err := repo.First(status)
	assert.NoError(t, err)
	assert.Equal(t, entity.Bad, status.Morning.Fitness)
}

func TestSave(t *testing.T) {
	repo := initTest(t)

	var err error
	var status, statusWithMorning, updatedStatusWithMorning, responseStatus *entity.Status
	var id1, id2 uint
	status = &entity.Status{}
	err = repo.Save(status)
	id1 = status.ID
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id1)

	statusWithMorning = &entity.Status{Morning: &entity.MorningStatus{Fitness: entity.Good}}
	err = repo.Save(statusWithMorning)
	id2 = statusWithMorning.ID
	assert.NoError(t, err)
	assert.EqualValues(t, 2, id2)

	updatedStatusWithMorning = &entity.Status{ID: id2, Morning: &entity.MorningStatus{ID: statusWithMorning.Morning.ID, Fitness: entity.VeryGood}}
	err = repo.Save(updatedStatusWithMorning)
	assert.NoError(t, err)
	assert.Equal(t, updatedStatusWithMorning.ID, uint(2))

	responseStatus = &entity.Status{ID: id2}
	repo.db.Preload(clause.Associations).First(responseStatus)
	assert.NotNil(t, responseStatus.Morning)
	assert.EqualValues(t, entity.VeryGood, responseStatus.Morning.Fitness)
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
