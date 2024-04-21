package api

import (
	"errors"

	"encore.app/internalAuth"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}

var histaDB *sqldb.Database = sqldb.NewDatabase("blog", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func initDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: histaDB.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	db, err := initDb()
	if err != nil {
		return nil, err
	}
	internalAuth.CurrentToken = internalAuth.CurrentToken.InitToken()
	if internalAuth.CurrentToken == nil {
		return nil, errors.New("Token not initialized")
	}
	internalAuth.InitUsers()

	return &Service{db: db}, nil
}
