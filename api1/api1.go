package api1

import (
	"context"

	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}

type Response struct {
	Resp string
}

//encore:api public method=GET path=/test
func Get(ctx context.Context) (Response, error) {
	return Response{Resp: "Hallo Julia, ich bin eine API, aber noch keine Datenbank"}, nil
}
