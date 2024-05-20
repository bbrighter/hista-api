package notes

import (
	"time"

	"encore.app/errors"
	"gorm.io/gorm/clause"
)

type Note struct {
	ID   uint
	Date time.Time
	Text string
}

type Notes []Note

// Create a new note with date = now
func createNote(service *Service) Note {
	var note = Note{
		Date: time.Now(),
	}
	service.db.Create(&note)
	return note
}

// Get all notes including text
func getNotes(service *Service) (Notes, error) {
	var notes = []Note{}
	var err error = service.db.Find(&notes).Error
	return notes, err
}

// Patch note. Date and text are optional
func (note Note) patch(service *Service, date *time.Time, text *string) (Note, error) {
	if note.ID == 0 {
		return Note{}, errors.ErrorIDMissing
	}
	tx := service.db.Model(&note)
	var newNote = note
	if date != nil {
		newNote.Date = *date
		tx.Select("date")
	}
	if text != nil {
		newNote.Text = *text
		tx.Select("text") // Make sure that text gets updated, even if ""
	}
	tx = tx.Clauses(clause.Returning{}).Updates(&newNote)
	if tx.RowsAffected == 0 {
		return Note{}, errors.ErrorNotFound
	}
	return note, tx.Error
}

// Delete a note
func (note Note) delete(service *Service) error {
	if note.ID == 0 {
		return errors.ErrorIDMissing
	}
	tx := service.db.Delete(&note)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}
