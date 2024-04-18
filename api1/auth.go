package api1

import (
	"context"
	"math/rand"
	"strings"

	"encore.dev/beta/errs"
)

type AuthParams struct {
	Password string
}

type Token struct {
	Token *string
}

var currentToken Token

//encore:api public method=POST path=/auth
func (s *Service) Auth(ctx context.Context, params AuthParams) (Token, error) {
	pw := params.Password
	if pw == "123pi" {
		sb := strings.Builder{}
		sb.Grow(50)
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		for i := 0; i < 50; i++ {
			sb.WriteByte(charset[rand.Intn(len(charset))])
		}
		var tok string = sb.String()
		currentToken.Token = &tok
		return currentToken, nil
	}
	return currentToken, &errs.Error{Code: errs.Unauthenticated}
}
