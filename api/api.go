package api

import (
	"context"

	"encore.dev/beta/errs"
)

type TodoItem struct {
	ID    int64  `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TodoItems struct {
	Items []TodoItem `json:"items"`
}

//encore:api auth method=GET path=/test
func (s *Service) Get(ctx context.Context) (TodoItems, error) {
	var todoItems []TodoItem
	err := s.db.Find(&todoItems).Error
	return TodoItems{todoItems}, err
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
