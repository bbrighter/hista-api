package notes

import (
	"context"
	"time"
)

type NotesResp struct {
	Notes []NoteResp `json:"notes"`
}

type NoteResp struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
}

// encore:api auth method=GET path=/notes
func (service *Service) GetNotes(ctx context.Context) (NotesResp, error) {
	notes := getNotes(service)
	return notes.toResp(), nil
}

// encore:api auth method=POST path=/notes
func (service *Service) PostNote(ctx context.Context) (NoteResp, error) {
	note := createNote(service)
	return note.toResp(), nil
}

// encore:api auth method=DELETE path=/notes/:noteId
func (service *Service) DeleteNote(ctx context.Context, noteId uint) error {
	var note = Note{ID: noteId}
	return note.delete(service)
}

type NoteParams struct {
	Date *time.Time `json:"date" encore:"optional"`
	Text *string    `json:"text" encore:"optional"`
}

// encore:api auth method=PATCH path=/notes/:noteId
func (service *Service) PatchNote(ctx context.Context, noteId uint, params NoteParams) (NoteResp, error) {
	var note, respNote Note
	note = Note{ID: noteId}

	var err error
	respNote, err = note.patch(service, params.Date, params.Text)
	return respNote.toResp(), err
}
