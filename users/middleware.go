package users

import (
	"encore.app/authentication"
	"encore.app/errors"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
)

const USER_MANAGEMENT_APP string = "user-management"

// encore:middleware target=tag:user
func ValidationMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	data, ok := auth.Data().(authentication.AuthData)
	if !ok {
		return next(req)
	}
	if !data.AppMapping[USER_MANAGEMENT_APP] {
		return middleware.Response{Err: errors.ErrorUnauthenticated}
	}
	return next(req)
}
