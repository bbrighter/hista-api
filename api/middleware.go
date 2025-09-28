package api

import (
	"context"

	"encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
)

// encore:middleware target=all
func AddPiidMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	data, ok := auth.Data().(*entity.AuthData)
	if !ok {
		return middleware.Response{Err: errors.ErrorUnauthenticated}
	}
	ctx := context.WithValue(req.Context(), "piid", data.PIID)
	reqWithCtx := req.WithContext(ctx)
	return next(reqWithCtx)
}
