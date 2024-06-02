package notes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNotesAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var err error
	var resp NotesResp
	resp, err = service.GetNotes(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Notes, 0)

	var note Note
	note.createNote(service)
	resp, err = service.GetNotes(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Notes, 1)
}

func TestDeleteNoteAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var err error
	var note Note
	note.createNote(service)

	err = service.DeleteNote(ctx, note.ID)
	assert.NoError(t, err)

	err = service.DeleteNote(ctx, 10000)
	assert.Error(t, err)
}

func TestPatchNoteAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var params NoteParams
	var err error
	var resp NoteResp
	var note Note
	note.createNote(service)

	resp, err = service.PatchNote(ctx, note.ID, params)
	assert.NoError(t, err)
	assert.Equal(t, note.Text, resp.Text)
	assert.True(t, note.Date.Truncate(time.Second).Equal(resp.Date.Truncate(time.Second)))

	newTime := time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local)
	params.Date = &newTime
	resp, err = service.PatchNote(ctx, note.ID, params)
	assert.NoError(t, err)
	assert.True(t, newTime.Equal(resp.Date))

	newText := "text"
	params.Date = nil
	params.Text = &newText
	resp, err = service.PatchNote(ctx, note.ID, params)
	assert.NoError(t, err)
	assert.Equal(t, "text", resp.Text)
}
