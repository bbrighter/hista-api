package product_mgmt

import (
	"encore.app/product_mgmt/product"
	"encore.app/product_mgmt/product_instance"
	"encore.dev/config"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	prod *product.ProductService
	pi   *product_instance.ProductInstanceService
}

var prodMgmtDb = sqldb.NewDatabase("product_mgmt_db", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

var cfg *product.Config = config.Load[*product.Config]()

func initDB() (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: prodMgmtDb.Stdlib(),
	}))
}

func initService() (*Service, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	pi := product_instance.NewProductInstanceService(db, cfg.ToProducts(), cfg.ToApps())
	prod := product.NewProductService(product.NewProductRepo(cfg.ToProducts(), cfg.ToApps()))

	var service = &Service{
		pi:   pi,
		prod: prod,
	}
	return service, nil
}
