package symptoms

import (
	"context"
	_ "embed"
	"log"

	"encore.dev"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type Service struct {
	db *gorm.DB
}

var histaDB = sqldb.Named("hista_db")

//go:embed fixtures.sql
var fixtures string

func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: histaDB.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	if encore.Meta().Environment.Cloud == encore.CloudLocal {
		if _, err := histaDB.Exec(context.Background(), fixtures); err != nil {
			log.Fatalln("unable to add fixtures:", err)
		}
	}

	return &Service{db: db}, nil
}
