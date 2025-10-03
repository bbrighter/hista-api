package notes

import (
	"context"
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"encore.app/generic_queries"
	"gorm.io/gorm/clause"
)

// Create a new note with date = now
func (repo *NotesRepository) Create(ctx context.Context, note *entity.Note) error {
	return generic_queries.Create(ctx, repo.db, note)
}

// Get all notes including text
func (repo *NotesRepository) List(ctx context.Context) ([]*entity.Note, error) {
	return generic_queries.List[*entity.Note](ctx, repo.db)
}

// Patch note. Date and text are optional
func (repo *NotesRepository) Patch(ctx context.Context, id uint, date *time.Time, text *string) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}

	var note = entity.Note{ID: id}
	tx := repo.db.Model(&note)
	if date == nil && text == nil {
		return errors.ErrorAttributeMustBeSet("date or text")
	}
	if date != nil {
		note.Date = *date
		tx.Select("date")
	}
	if text != nil {
		note.Text = *text
		tx.Select("text") // Make sure that text gets updated, even if ""
	}
	tx = tx.Clauses(clause.Returning{}).Where("pi_id = ?", piid).Updates(&note)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

// Delete a note
func (repo *NotesRepository) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Note](ctx, repo.db, id)
}
