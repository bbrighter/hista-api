package api1

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}

// Migrations disabled, get DB to run
var blogDB = sqldb.NewDatabase("blog", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: blogDB.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	return &Service{db: db}, nil
}

type Version struct {
	Resp    string
	Version string
}

//encore:api public method=GET path=/test
func (s *Service) Get(ctx context.Context) (TodoItems, error) {
	var todoItems []TodoItem
	err := s.db.Find(&todoItems).Error
	return TodoItems{todoItems}, err
}

type TodoItem struct {
	ID    int64  `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TodoItems struct {
	Items []TodoItem `json:"items"`
}

type PostParams struct {
	Title string `json:"title"`
}

//encore:api auth method=POST path=/test
func (s *Service) Post(ctx context.Context, params PostParams) error {
	if params.Title == "" {
		return &errs.Error{Code: errs.InvalidArgument, Message: "Name needed"}
	}
	var todoItem = TodoItem{Title: params.Title, Done: false}
	return s.db.Create(&todoItem).Error
}

// AuthHandler can be named whatever you prefer (but must be exported).
//
//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, error) {
	if token != *currentToken.Token {
		return "", &errs.Error{Code: errs.Unauthenticated}
	}
	return "julia", nil
}
