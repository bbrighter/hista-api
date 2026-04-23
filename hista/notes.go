package hista

import (
	"context"
	"sort"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/notes"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type NoteListResp struct {
	Notes []NoteResp `json:"notes"`
}

type NoteResp struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
}

func toNoteListResp(notes []*notes.Note) NoteListResp {
	var resps = []NoteResp{}
	for _, note := range notes {
		var resp NoteResp = toNoteResp(note)
		resps = append(resps, resp)
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return NoteListResp{Notes: resps}
}

func toNoteResp(note *notes.Note) NoteResp {
	return NoteResp{
		ID:   note.ID,
		Date: note.Date,
		Text: note.Text,
	}
}

// encore:api auth method=GET path=/piid/:piid/notes
func (service *Service) ListNotes(ctx context.Context, piid uuid.UUID) (NoteListResp, error) {
	notes, err := service.notes.ListNotes(ctx)
	if err != nil {
		return NoteListResp{}, errors.MapError(err)
	}
	return toNoteListResp(notes), nil
}

// encore:api auth method=POST path=/piid/:piid/notes
func (service *Service) PostNote(ctx context.Context, piid uuid.UUID) (NoteResp, error) {
	date := time.Now()
	id, err := service.notes.CreateNote(ctx, date)
	if err != nil {
		return NoteResp{}, errors.MapError(err)
	}
	return toNoteResp(&notes.Note{ID: id, Date: date}), nil
}

// encore:api auth method=DELETE path=/piid/:piid/notes/:noteId
func (service *Service) DeleteNote(ctx context.Context, piid uuid.UUID, noteId uint) error {
	return errors.MapError(service.notes.DeleteNote(ctx, noteId))
}

type NoteParams struct {
	Date option.Option[time.Time] `json:"date" encore:"optional"`
	Text option.Option[string]    `json:"text" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/notes/:noteId
func (service *Service) PatchNote(ctx context.Context, piid uuid.UUID, noteId uint, params NoteParams) error {
	if params.Date.IsNone() && params.Text.IsNone() {
		return errors.BadRequest("at least one of 'date' or 'text' must be set")
	}
	err := service.notes.UpdateNote(ctx, noteId, params.Date.PtrOrNil(), params.Text.PtrOrNil())
	return errors.MapError(err)
}
