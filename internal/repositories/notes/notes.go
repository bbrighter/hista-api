package notes

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

// Create a new note with date = now
func (repo *NotesRepository) Create() (uint, error) {
	var note = entity.Note{Date: time.Now()}
	err := repo.db.Create(&note).Error
	return note.ID, err
}

// Get all notes including text
func (repo *NotesRepository) List() entity.Notes {
	var notes = entity.Notes{}
	repo.db.Find(&notes)
	return notes
}

// Patch note. Date and text are optional
func (repo *NotesRepository) Patch(id uint, date *time.Time, text *string) error {
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
	tx = tx.Clauses(clause.Returning{}).Updates(&note)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

// Delete a note
func (repo *NotesRepository) Delete(id uint) error {
	tx := repo.db.Delete(&entity.Note{ID: id})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}
