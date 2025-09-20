package authentication

import (
	"context"

	"encore.app/errors"
	"encore.dev/beta/auth"
	"encore.dev/types/uuid"
)

type AuthParams struct {
	token string `header:"Auth"`
}

type AuthData struct {
	AppMapping map[string]bool
	PIID       uuid.UUID
}

//encore:authhandler
func AuthHandler(ctx context.Context, params *AuthParams) (auth.UID, *AuthData, error) {
	token, err := parseToken(params.token, keys.publicKey)
	if err != nil {
		return auth.UID(""), &AuthData{}, errors.ErrorUnauthenticated
	}
	var appMapping = make(map[string]bool)

	authData := &AuthData{
		AppMapping: appMapping,
		PIID:       uuid.UUID{},
	}

	return auth.UID(token.Subject), authData, nil
}
