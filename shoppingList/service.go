package shoppinglist

import (
	"time"

	"encore.app/shoppingList/internal"
	"encore.app/shoppingList/internal/repository"
	unitofwork "encore.app/shoppingList/internal/unitOfWork"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	item internal.IItemUseCase
	list internal.IListUseCase
	prod internal.IProductUseCase
}

var shoppingListDb *sqldb.Database = sqldb.NewDatabase("shopping_list", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func initDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: shoppingListDb.Stdlib(),
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
	// unitOfWork := repository.NewUnitOfWork(db)
	uow := unitofwork.NewUnitOfWork(db)
	itemRepo := repository.NewItemRepo(db)
	productRepo := repository.NewProductRepo(db)
	listRepo := repository.NewListRepo(db)
	uc := internal.NewItemUseCase(itemRepo, productRepo, uow)
	list := internal.NewListUseCase(listRepo, itemRepo, uow)
	products := internal.NewProductUseCase(productRepo)
	return &Service{item: uc, list: list, prod: products}
}
