package internalAuth

import (
	"context"

	"encore.app/errors"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

type AuthParams struct {
	Authorization string `header:"Authorization"`
	UserName      string `header:"UserName"`
}

//encore:authhandler
func AuthHandler(ctx context.Context, params *AuthParams) (auth.UID, error) {
	var uid uuid.UUID
	var err error
	uid, err = uuid.FromString(params.Authorization)
	if err != nil {
		return "", errors.NewError("invalid uuid", errs.InvalidArgument)
	}
	user, err := getUserByName(params.UserName)
	if err != nil {
		return "", errors.ErrorUnauthenticated
	}
	var token = Token{
		Bearer: uid,
		UserID: user.ID,
	}
	var valid bool
	valid, err = token.isValid()
	if !valid {
		return "", errors.ErrorNotFound
	}
	return auth.UID(user.Name), err
}
