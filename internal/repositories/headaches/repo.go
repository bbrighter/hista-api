package headaches

import "gorm.io/gorm"

type HeadacheRepository struct {
	db *gorm.DB
}

func NewHeadacheRepository(db *gorm.DB) *HeadacheRepository {
	return &HeadacheRepository{db: db}
}
