package api

import (
	"time"

	"encore.app/api/repositories/meals"
	"encore.app/api/repositories/notes"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	DB    *gorm.DB
	meal  *meals.MealRepository
	notes *notes.NotesRepository
}

var HistaDB *sqldb.Database = sqldb.NewDatabase("hista_db", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func initDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: HistaDB.Stdlib(),
	}))

	if err != nil {
		return nil, err
	}
	sqlDb, _ := db.DB()
	sqlDb.SetConnMaxLifetime(time.Second)
	return db, nil
}

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	db, err := initDb()
	if err != nil {
		return nil, err
	}
	return &Service{
		DB:    db,
		meal:  meals.NewMealRepository(db),
		notes: notes.NewNotesRepository(db),
	}, nil
}
