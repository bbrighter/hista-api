package notes

import "gorm.io/gorm"

type NotesRepository struct {
	db *gorm.DB
}

func NewNotesRepository(db *gorm.DB) *NotesRepository {
	return &NotesRepository{db: db}
}
