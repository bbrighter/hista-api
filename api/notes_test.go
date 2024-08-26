package api

import (
	"context"
	"testing"
	"time"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func (service *Service) createTestNote(t *testing.T) (uint, func()) {
	ctx := context.TODO()
	note, err := service.PostNote(ctx)
	assert.NoError(t, err)
	cleanUp := func() {
		service.DeleteNote(ctx, note.ID)
	}
	return note.ID, cleanUp
}

func TestGetNotes(t *testing.T) {
	service, ctx := initAPITest(t)

	var resp entity.NotesResp
	resp, _ = service.GetNotes(ctx)
	assert.Len(t, resp.Notes, 0)

	_, cleanUp := service.createTestNote(t)
	defer cleanUp()
	resp, _ = service.GetNotes(ctx)
	assert.Len(t, resp.Notes, 1)
}

func TestDeleteNoteAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteNote(ctx, 1)
	assert.Error(t, err)

	id, cleanUp := service.createTestNote(t)
	defer cleanUp()

	err = service.DeleteNote(ctx, id)
	assert.NoError(t, err)
}

func TestPatchNoteAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params NoteParams
	var err error
	err = service.PatchNote(ctx, 1, params)
	assert.Error(t, err)

	id, cleanup := service.createTestNote(t)
	defer cleanup()

	err = service.PatchNote(ctx, id, params)
	assert.Error(t, err)

	newTime := time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local)
	params.Date = &newTime
	err = service.PatchNote(ctx, id, params)
	assert.NoError(t, err)

	newText := "text"
	params.Date = nil
	params.Text = &newText
	err = service.PatchNote(ctx, id, params)
	assert.NoError(t, err)
}
