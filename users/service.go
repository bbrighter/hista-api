package users

import (
	"encore.app/users/internal"
	"encore.app/users/internal/repository"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	db      *gorm.DB
	user    internal.IUser
	mgmt    internal.IUserManagement
	auth    internal.IAuth
	secrets *Secrets
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
	r := repository.NewUserRepo(db, cost)
	user := internal.NewUserUseCase(r)
	mgmt := internal.NewUserManagement(r)
	auth := internal.NewAuthUseCase(r)
	secrets := NewSecrets()

	var service = &Service{
		db:      db,
		user:    user,
		mgmt:    mgmt,
		auth:    auth,
		secrets: secrets,
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
