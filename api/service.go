package api

import (
	"time"

	"encore.app/internal"
	"encore.app/internal/repositories/meals"
	"encore.app/internal/repositories/notes"
	"encore.app/internal/repositories/symptoms"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	DB              *gorm.DB
	mealUC          internal.IMealUseCase
	food            internal.IFoodUseCase
	notes           *notes.NotesRepository
	symtpoms        *symptoms.SymptomsRepo
	conditionEvents internal.IConditionEventUseCase
	conditions      internal.IConditionUseCase
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
	mealRepo := meals.NewMealRepository(db)
	symptomRepo := symptoms.NewSymtpomsRepo(db)

	return &Service{
		DB:              db,
		mealUC:          internal.NewMealUseCase(mealRepo),
		food:            internal.NewFoodUseCase(mealRepo, mealRepo),
		notes:           notes.NewNotesRepository(db),
		symtpoms:        symptoms.NewSymtpomsRepo(db),
		conditionEvents: internal.NewConditionEventUseCase(symptomRepo),
		conditions:      internal.NewconditionsUseCase(symptomRepo, symptomRepo),
	}, nil
}
