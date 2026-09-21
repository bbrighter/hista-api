package users

import (
	userproductinstance "encore.app/users/internal/user_product_instance"
	"encore.app/users/internal/users"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	u   *users.UserService
	upi *userproductinstance.UserProductInstanceService
}

var usersDB = sqldb.NewDatabase("users_db", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func initDB() (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: usersDB.Stdlib(),
	}))
}

func setupService(db *gorm.DB, cost int) *Service {

	u := users.NewUserService(db, cost)

	r_upi := userproductinstance.NewUserProductInstanceRepo(db)
	upi := userproductinstance.NewUserProductInstanceService(r_upi)

	var service = &Service{
		u:   u,
		upi: upi,
	}
	return service
}

func initService() (*Service, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}
	service := setupService(db, 15)

	return service, nil
}
