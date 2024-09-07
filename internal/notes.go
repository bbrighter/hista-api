package internal

import (
	"time"

	"encore.app/entity"
)

type NoteUseCase struct {
	repo INotesRepository
}

func NewNoteUseCase(repo INotesRepository) NoteUseCase {
	return NoteUseCase{repo: repo}
}

func (uc NoteUseCase) List() entity.Notes {
	return uc.repo.List()
}

func (uc NoteUseCase) Create() (entity.Note, error) {
	return uc.repo.Create()
}

func (uc NoteUseCase) Patch(id uint, date *time.Time, text *string) error {
	return uc.repo.Patch(id, date, text)
}

func (uc NoteUseCase) Delete(id uint) error {
	return uc.repo.Delete(id)
}
