package symptoms

import "gorm.io/gorm"

type SymptomsRepo struct {
	db *gorm.DB
}

func NewSymtpomsRepo(db *gorm.DB) *SymptomsRepo {
	return &SymptomsRepo{db: db}
}
