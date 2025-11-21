package internal

import (
	"context"
	"time"

	"encore.app/hista/entity"
)

type (
	INotesRepository interface {
		List(ctx context.Context) ([]*entity.Note, error)
		Create(ctx context.Context, note *entity.Note) error
		Patch(ctx context.Context, id uint, date *time.Time, text *string) error
		Delete(ctx context.Context, id uint) error
	}

	INotesUseCase interface {
		List(ctx context.Context) (entity.Notes, error)
		Create(ctx context.Context) (entity.Note, error)
		Patch(ctx context.Context, id uint, date *time.Time, text *string) error
		Delete(ctx context.Context, id uint) error
	}
)

type NoteUseCase struct {
	repo INotesRepository
}

func NewNoteUseCase(repo INotesRepository) NoteUseCase {
	return NoteUseCase{repo: repo}
}

func (uc NoteUseCase) List(ctx context.Context) (entity.Notes, error) {
	notes, err := uc.repo.List(ctx)
	return notes, err
}

func (uc NoteUseCase) Create(ctx context.Context) (entity.Note, error) {
	var note = &entity.Note{Date: time.Now(), Text: ""}
	err := uc.repo.Create(ctx, note)
	return *note, err
}

func (uc NoteUseCase) Patch(ctx context.Context, id uint, date *time.Time, text *string) error {
	return uc.repo.Patch(ctx, id, date, text)
}

func (uc NoteUseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
