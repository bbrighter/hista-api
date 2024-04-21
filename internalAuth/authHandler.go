package internalAuth

import (
	"context"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

//encore:authhandler
func AuthHandler(ctx context.Context, bearer string) (auth.UID, error) {
	token := Token{
		UserId:  Julia,
		Bearer:  bearer,
		Expires: time.Now(),
	}
	if token.isValidToken() {
		return "Julia", nil
	} else {
		return "", &errs.Error{Code: errs.Unauthenticated}
	}
}
