package users

import (
	"encore.app/users/users"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	db *gorm.DB
	u  *users.UserService
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

	var service = &Service{
		u: u,
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
