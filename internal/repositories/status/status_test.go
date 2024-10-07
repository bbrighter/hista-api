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
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Status{})
	assert.NoError(t, err)
	return NewStatusRepo(db)
}

func TestFind(t *testing.T) {
	repo := initTest(t)

	var statuses []entity.Status
	statuses = repo.Find()
	assert.Len(t, statuses, 0)

	status := entity.Status{Date: time.Now(), TimeOfDay: entity.Evening, Fitness: entity.Bad}
	repo.db.Create(&status)
	statuses = repo.Find()
	assert.Len(t, statuses, 1)
}

func TestCreate(t *testing.T) {
	repo := initTest(t)

	status := &entity.Status{}
	err := repo.Create(status)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, status.ID)
}
