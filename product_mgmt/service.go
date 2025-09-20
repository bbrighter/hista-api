package product_mgmt

import (
	"encore.app/product_mgmt/internal"
	"encore.app/product_mgmt/internal/repository"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	instance internal.InstanceStore
	product  internal.ProductFinder
}

var usersDB = sqldb.NewDatabase("product_mgmt_db", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func initDB() (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: usersDB.Stdlib(),
	}))
}

func initService() (*Service, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	config := repository.ParseConfig("config.yaml")
	p := repository.NewProductRepo(config)
	i := repository.NewInstanceRepo(db)

	instance := internal.NewInstanceUseCase(i, p)
	product := internal.NewProductFinder(p)

	var service = &Service{
		instance: instance,
		product:  product,
	}
	return service, nil
}
