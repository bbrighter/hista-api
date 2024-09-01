package api

import (
	"time"

	"encore.app/internal"
	"encore.app/internal/repositories/meals"
	"encore.app/internal/repositories/notes"
	"encore.app/internal/repositories/statistics"
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
	notes           internal.INotesUseCase
	symtpoms        internal.ISymptomsUseCase
	conditionEvents internal.IConditionEventUseCase
	conditions      internal.IConditionUseCase
	diary           internal.IDiaryUseCase
	statistics      internal.IStatisticsUseCase
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
	noteRepo := notes.NewNotesRepository(db)
	statsRepo := statistics.NewStatisticsRepo(db)

	return &Service{
		DB:              db,
		mealUC:          internal.NewMealUseCase(mealRepo),
		food:            internal.NewFoodUseCase(mealRepo, mealRepo),
		notes:           internal.NewNoteUseCase(noteRepo),
		symtpoms:        internal.NewSymptomsUseCase(symptomRepo, symptomRepo),
		conditionEvents: internal.NewConditionEventUseCase(symptomRepo),
		conditions:      internal.NewconditionsUseCase(symptomRepo, symptomRepo),
		diary:           internal.NewDiaryUseCase(mealRepo, symptomRepo, symptomRepo, noteRepo),
		statistics:      internal.NewStatisticsUseCase(statsRepo),
	}, nil
}
