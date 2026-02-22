package hista

import (
	"time"

	"encore.app/hista/internal"
	"encore.app/hista/internal/repositories/headaches"
	"encore.app/hista/internal/repositories/meals"
	"encore.app/hista/internal/repositories/medicine"
	"encore.app/hista/internal/repositories/move"
	"encore.app/hista/internal/repositories/notes"
	"encore.app/hista/internal/repositories/pollen"
	"encore.app/hista/internal/repositories/statistics"
	"encore.app/hista/internal/repositories/status"
	"encore.app/hista/internal/repositories/symptoms"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	DB                 *gorm.DB
	meals              internal.IMealUseCase
	foods              internal.IFoodUseCase
	ingredients        internal.IIngredientUseCase
	ingredientsManager internal.IIngredientManager
	notes              internal.INotesUseCase
	symptoms           internal.ISymptomsUseCase
	conditionEvents    internal.IConditionEventUseCase
	conditions         internal.IConditionUseCase
	diary              internal.IDiaryUseCase
	statistics         internal.IStatisticsUseCase
	pollens            internal.IPollenUseCase
	status             internal.IStatusUseCase
	headaches          internal.IHeadacheUseCase
	move               internal.PiidMover
	medicineList       internal.MedicineLister
	medicine           internal.MedicineManager
	intake             internal.IntakeManager
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
	return initServiceWithDb(db), nil
}

func initServiceWithDb(db *gorm.DB) *Service {
	mealRepo := meals.NewMealRepository(db)
	symptomRepo := symptoms.NewSymptomsRepo(db)
	noteRepo := notes.NewNotesRepository(db)
	statsRepo := statistics.NewStatisticsRepo(db)
	pollenRepo := pollen.NewPollenRepo(db)
	dwdRepo := pollen.NewDWDRepo()
	statusRepo := status.NewStatusRepo(db)
	headacheRepo := headaches.NewHeadacheRepository(db)
	moveRepo := move.NewMoveRepo(db)
	medicineRepo := medicine.NewMedicineRepo(db)
	intakeRepo := medicine.NewIntakeRepo(db)

	return &Service{
		DB:                 db,
		meals:              internal.NewMealUseCase(mealRepo, mealRepo),
		ingredients:        internal.NewIngredientUseCase(mealRepo),
		ingredientsManager: internal.NewIngredientsManager(mealRepo),
		foods:              internal.NewFoodUseCase(mealRepo, mealRepo),
		notes:              internal.NewNoteUseCase(noteRepo),
		symptoms:           internal.NewSymptomsUseCase(symptomRepo, symptomRepo),
		conditionEvents:    internal.NewConditionEventUseCase(symptomRepo, symptomRepo),
		conditions:         internal.NewConditionsUseCase(symptomRepo, symptomRepo),
		diary:              internal.NewDiaryUseCase(mealRepo, symptomRepo, symptomRepo, noteRepo, pollenRepo, intakeRepo),
		statistics:         internal.NewStatisticsUseCase(statsRepo),
		pollens:            internal.NewPollenUseCase(pollenRepo, dwdRepo),
		status:             internal.NewStatusUseCase(statusRepo),
		headaches:          internal.NewHeadacheUseCase(headacheRepo),
		move:               internal.NewPiidMoveUseCase(moveRepo),
		medicineList:       internal.NewMedicineListUseCase(medicineRepo),
		medicine:           internal.NewMedicineMgtmUseCase(medicineRepo),
		intake:             internal.NewIntakeMgmtUseCase(intakeRepo),
	}
}
