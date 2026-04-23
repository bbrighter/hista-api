package notes

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type NotesRepo struct {
	db *gorm.DB
}

func NewNotesRepo(db *gorm.DB) *NotesRepo {
	return &NotesRepo{db: db}
}

func (r *NotesRepo) ListNotes(ctx context.Context) ([]*Note, error) {
	return generic_queries.List[*Note](ctx, r.db)
}

func (r *NotesRepo) CreateNote(ctx context.Context, n *Note) error {
	return generic_queries.Create(ctx, r.db, n)
}

func (r *NotesRepo) DeleteNote(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Note](ctx, r.db, id)
}

func (r *NotesRepo) UpdateNote(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "notes", id, values)
}
