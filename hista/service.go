package hista

import (
	"time"

	"encore.app/hista/internal/diary"
	"encore.app/hista/internal/dwd"
	"encore.app/hista/internal/dwdPollen"
	"encore.app/hista/internal/headaches"
	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/statistics"
	"encore.app/hista/internal/status"
	"encore.app/hista/internal/symptoms"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//encore:service
type Service struct {
	DB *gorm.DB

	meals     *meals.MealService
	stats     *statistics.StatisticsService
	meds      *medicines.MedicinesService
	syms      *symptoms.SymptomService
	notes     *notes.NotesService
	headaches *headaches.HeadacheService
	status    *status.StatusService
	pollens   *dwdPollen.DWDPollenService
	diary     *diary.DiaryService
}

var HistaDB *sqldb.Database = sqldb.NewDatabase("hista_db", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func initDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: HistaDB.Stdlib(),
	}), &gorm.Config{
		TranslateError: true,
		Logger:         EncoreLogger{level: logger.Silent},
	})

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
	dwdClient := dwd.NewDWDClient()
	return initServiceWithDb(db, dwdClient), nil
}

func initServiceWithDb(db *gorm.DB, dwdClient dwdPollen.KarlsruheDataGetter) *Service {
	medRepo := medicines.NewMedicineRepo(db)
	mealRepo := meals.NewMealRepository(db)
	symptomRepo := symptoms.NewSymptomRepo(db)
	condRepo := symptoms.NewConditionRepo(db)
	notesRepo := notes.NewNotesRepo(db)
	pollenRepo := pollen.NewPollenRepo(db)

	return &Service{
		DB: db,

		meals:     meals.NewMealService(db),
		stats:     statistics.NewStatisticsService(db),
		meds:      medicines.NewMedicineService(db),
		syms:      symptoms.NewSymptomService(db),
		notes:     notes.NewNotesService(db),
		headaches: headaches.NewHeadacheService(db),
		status:    status.NewStatusService(db),
		pollens:   dwdPollen.NewDWDPollenService(db, dwdClient),
		diary:     diary.NewDiaryService(mealRepo, condRepo, symptomRepo, notesRepo, pollenRepo, medRepo),
	}
}
