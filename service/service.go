package service

import (
	"errors"

	"encore.app/internalAuth"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	DB *gorm.DB
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

	return &Service{DB: db}, nil
}
