package statistics

import "gorm.io/gorm"

type StatisticsRepo struct {
	db *gorm.DB
}

func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}
