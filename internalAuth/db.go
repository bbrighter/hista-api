package internalAuth

import (
	"errors"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	db *gorm.DB
}

var histaDB = sqldb.Named("hista_db")
var memorizedUsers []*User

func initService() (*Service, error) {
	var err error
	var db *gorm.DB
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: histaDB.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	var service = &Service{db: db}
	if err = service.initializeUsers(); err != nil {
		return nil, err
	}
	updateMemorizedUsers(service)
	if len(memorizedUsers) == 0 {
		return nil, errors.New("No users")
	}
	return service, nil
}
