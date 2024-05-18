package notes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var note Note = createNote(service)
	assert.True(t, note.Date.Before(time.Now()))
}

func TestGetNotes(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var notes []Note
	notes = getNotes(service)
	assert.Len(t, notes, 0)

	createNote(service)
	notes = getNotes(service)
	assert.Len(t, notes, 1)
}

func TestPatchNote(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var note, newNote, nonExistingNote, noIDNote Note
	var err error
	note = createNote(service)
	newNote, err = note.patch(service, nil, nil)
	assert.NoError(t, err)
	assert.True(t, note.Date.Truncate(time.Second).Equal(newNote.Date.Truncate(time.Second)))

	newTime := time.Now().Add(-time.Hour)
	newNote, err = note.patch(service, &newTime, nil)
	assert.NoError(t, err)
	assert.True(t, note.Date.After(newNote.Date))

	newContent := "text"
	newNote, err = note.patch(service, nil, &newContent)
	assert.NoError(t, err)
	assert.Equal(t, newContent, newNote.Text)

	newContent2 := "text2"
	newTime2 := time.Now().Add(-time.Hour * 2)
	newNote, err = note.patch(service, &newTime2, &newContent2)
	assert.NoError(t, err)
	assert.Equal(t, newContent2, newNote.Text)
	assert.True(t, note.Date.After(newNote.Date))

	emptyContent := ""
	newNote, err = note.patch(service, nil, &emptyContent)
	assert.NoError(t, err)
	assert.Equal(t, emptyContent, newNote.Text)

	nonExistingNote = Note{ID: 100000}
	_, err = nonExistingNote.patch(service, nil, nil)
	assert.Error(t, err)

	noIDNote = Note{}
	_, err = noIDNote.patch(service, nil, nil)
	assert.Error(t, err)
}

func TestDeleteNote(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var note, noIDNote Note
	note = createNote(service)

	var err error
	err = note.delete(service)
	assert.NoError(t, err)

	err = note.delete(service)
	assert.Error(t, err)

	noIDNote = Note{}
	err = noIDNote.delete(service)
	assert.Error(t, err)
}
