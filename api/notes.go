package api

import (
	"context"
	"time"

	entity "encore.app/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/notes
func (service *Service) ListNotes(ctx context.Context, piid uuid.UUID) (entity.NotesResp, error) {
	notes, err := service.notes.List(ctx)
	return notes.ToResp(), err
}

// encore:api auth method=POST path=/piid/:piid/notes
func (service *Service) PostNote(ctx context.Context, piid uuid.UUID) (entity.NoteResp, error) {
	note, err := service.notes.Create(ctx)
	if err != nil {
		return entity.NoteResp{}, err
	}
	return note.ToResp(), nil
}

// encore:api auth method=DELETE path=/piid/:piid/notes/:noteId
func (service *Service) DeleteNote(ctx context.Context, piid uuid.UUID, noteId uint) error {
	return service.notes.Delete(ctx, noteId)
}

type NoteParams struct {
	Date *time.Time `json:"date" encore:"optional"`
	Text *string    `json:"text" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/notes/:noteId
func (service *Service) PatchNote(ctx context.Context, piid uuid.UUID, noteId uint, params NoteParams) error {
	return service.notes.Patch(ctx, noteId, params.Date, params.Text)
}
