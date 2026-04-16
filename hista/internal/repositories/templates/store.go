package templates

import "gorm.io/gorm"

type TemplateStore struct {
	db *gorm.DB
}

func NewTemplateStore(db *gorm.DB) *TemplateStore {
	return &TemplateStore{db: db}
}
