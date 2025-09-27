package api

import (
	"context"
	"time"

	entity "encore.app/entity"
)

// encore:api auth method=GET path=/notes
func (service *Service) GetNotes(ctx context.Context) (entity.NotesResp, error) {
	notes, err := service.notes.List(ctx)
	return notes.ToResp(), err
}

// encore:api auth method=POST path=/notes
func (service *Service) PostNote(ctx context.Context) (entity.NoteResp, error) {
	note, err := service.notes.Create(ctx)
	if err != nil {
		return entity.NoteResp{}, err
	}
	return note.ToResp(), nil
}

// encore:api auth method=DELETE path=/notes/:noteId
func (service *Service) DeleteNote(ctx context.Context, noteId uint) error {
	return service.notes.Delete(ctx, noteId)
}

type NoteParams struct {
	Date *time.Time `json:"date" encore:"optional"`
	Text *string    `json:"text" encore:"optional"`
}

// encore:api auth method=PATCH path=/notes/:noteId
func (service *Service) PatchNote(ctx context.Context, noteId uint, params NoteParams) error {
	return service.notes.Patch(ctx, noteId, params.Date, params.Text)
}
