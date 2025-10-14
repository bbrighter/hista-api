package notes

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var GUID = uuid.FromStringOrNil("cf0d4408-8db5-4572-b5d9-4ed873d1341f")

func initTest(t *testing.T) (*NotesRepository, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Note{})
	assert.NoError(t, err)

	ctx := context.WithValue(t.Context(), contextKeys.Piid, GUID)
	return &NotesRepository{db: db}, ctx
}
func TestCreateNote(t *testing.T) {
	repo, ctx := initTest(t)

	now := time.Now()
	var note = &entity.Note{Date: now, Text: "Text"}
	err := repo.Create(ctx, note)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), note.ID)
	assert.Equal(t, "Text", note.Text)
	assert.True(t, now.Equal(note.Date))
}

func TestGetNotes(t *testing.T) {
	repo, ctx := initTest(t)

	notes, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, notes, 0)

	repo.db.Create(&entity.Note{ID: 1, PIID: GUID})
	notes, err = repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, notes, 1)
}

func TestPatchNote(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	err = repo.Patch(ctx, 1, nil, nil)
	assert.Error(t, err)

	repo.db.Create(&entity.Note{ID: 1, Date: time.Date(2000, 10, 9, 8, 7, 6, 0, time.Local), Text: "init", PIID: GUID})
	err = repo.Patch(ctx, 1, nil, nil)
	assert.Error(t, err)

	now := time.Now()
	err = repo.Patch(ctx, 1, &now, nil)
	assert.NoError(t, err)

	text := "text"
	err = repo.Patch(ctx, 1, nil, &text)
	assert.NoError(t, err)

	err = repo.Patch(ctx, 1, &now, &text)
	assert.NoError(t, err)
}

func TestDeleteNote(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	err = repo.Delete(ctx, 1)
	assert.Error(t, err)

	repo.db.Create(&entity.Note{ID: 1, PIID: GUID})
	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
}
