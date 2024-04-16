package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func InitRepository() Repository {
	return Repository{
		Db: initDb(),
	}
}

func initDb() *gorm.DB {
	// dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=prefer"
	dsn := "host=ep-patient-dream-a25t1r5n-pooler.eu-central-1.aws.neon.tech user=default password=yux50gdiAbUX dbname=verceldb port=5432 ssmode=require"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db
}
