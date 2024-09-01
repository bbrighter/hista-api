package notes

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTest(t *testing.T) *NotesRepository {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Note{})
	assert.NoError(t, err)
	return &NotesRepository{db: db}
}
func TestCreateNote(t *testing.T) {
	repo := initTest(t)

	id, err := repo.Create()
	assert.NoError(t, err)
	assert.Greater(t, id, uint(0))
}

func TestGetNotes(t *testing.T) {
	repo := initTest(t)

	var notes entity.Notes
	notes = repo.List()
	assert.Len(t, notes, 0)

	repo.db.Create(&entity.Note{ID: 1})
	notes = repo.List()
	assert.Len(t, notes, 1)
}

func TestPatchNote(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.Patch(1, nil, nil)
	assert.Error(t, err)

	repo.db.Create(&entity.Note{ID: 1, Date: time.Date(2000, 10, 9, 8, 7, 6, 0, time.Local), Text: "init"})
	err = repo.Patch(1, nil, nil)
	assert.Error(t, err)

	now := time.Now()
	err = repo.Patch(1, &now, nil)
	assert.NoError(t, err)

	text := "text"
	err = repo.Patch(1, nil, &text)
	assert.NoError(t, err)

	err = repo.Patch(1, &now, &text)
	assert.NoError(t, err)
}

func TestDeleteNote(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.Delete(1)
	assert.Error(t, err)

	repo.db.Create(&entity.Note{ID: 1})
	err = repo.Delete(1)
	assert.NoError(t, err)
}
