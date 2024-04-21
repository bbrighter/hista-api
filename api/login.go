package api

import (
	"context"

	"encore.app/internalAuth"
	"encore.dev/beta/errs"
)

//encore:api public method=POST path=/auth
func Login(ctx context.Context, params internalAuth.AuthParams) (*internalAuth.Token, error) {
	if internalAuth.IsValidCredential(params.Password, params.UserId) {
		internalAuth.CurrentToken.InitToken()
		return internalAuth.CurrentToken, nil
	}
	return internalAuth.CurrentToken, &errs.Error{Code: errs.Unauthenticated}
}
