package notes

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type NotesService struct {
	n *NotesRepo
}

func NewNotesService(db *gorm.DB) *NotesService {
	n := NewNotesRepo(db)
	return &NotesService{n: n}
}

func (s *NotesService) ListNotes(ctx context.Context) ([]*Note, error) {
	return s.n.ListNotes(ctx)
}

func (s *NotesService) CreateNote(ctx context.Context, date time.Time) (uint, error) {
	note := Note{Date: date}
	err := s.n.CreateNote(ctx, &note)
	return note.ID, err
}

func (s *NotesService) DeleteNote(ctx context.Context, id uint) error {
	return s.n.DeleteNote(ctx, id)
}

func (s *NotesService) UpdateNote(ctx context.Context, id uint, date *time.Time, text *string) error {
	values := make(map[string]any)
	if date != nil {
		values["date"] = *date
	}
	if text != nil {
		values["text"] = *text
	}
	return s.n.UpdateNote(ctx, id, values)

}
