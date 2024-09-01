package api

import (
	"context"
	"testing"
	"time"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
)

var testNote = new(entity.Note)

func (service *Service) createTestNote(t *testing.T) func(t *testing.T) {
	ctx := context.TODO()
	resp, err := service.PostNote(ctx)
	assert.NoError(t, err)
	testNote.ID = resp.ID
	cleanUp := func(t *testing.T) {
		err := service.DeleteNote(ctx, resp.ID)
		assert.NoError(t, err)
		testNote = new(entity.Note)
	}
	return cleanUp
}

func TestGetNotes(t *testing.T) {
	service, ctx := initAPITest(t)

	var resp entity.NotesResp
	resp, _ = service.GetNotes(ctx)
	assert.Len(t, resp.Notes, 0)

	cleanUp := service.createTestNote(t)
	defer cleanUp(t)
	resp, _ = service.GetNotes(ctx)
	assert.Len(t, resp.Notes, 1)
}

func TestDeleteNote(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteNote(ctx, 1)
	assert.Error(t, err)

	service.createTestNote(t)

	err = service.DeleteNote(ctx, testNote.ID)
	assert.NoError(t, err)
}

func TestPatchNote(t *testing.T) {
	service, ctx := initAPITest(t)

	var params NoteParams
	var err error
	err = service.PatchNote(ctx, 1, params)
	assert.Error(t, err)

	cleanup := service.createTestNote(t)
	defer cleanup(t)

	err = service.PatchNote(ctx, testNote.ID, params)
	assert.Error(t, err)

	newTime := time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local)
	params.Date = &newTime
	err = service.PatchNote(ctx, testNote.ID, params)
	assert.NoError(t, err)

	newText := "text"
	params.Date = nil
	params.Text = &newText
	err = service.PatchNote(ctx, testNote.ID, params)
	assert.NoError(t, err)
}
