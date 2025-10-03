package api

import (
	"context"

	entity "encore.app/entity"
	"encore.app/errors"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
	uuid "encore.dev/types/uuid"
)

// encore:middleware target=tag:external
func AddPiidMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	params := req.Data().PathParams
	piidStr := params.Get("piid")
	url_piid, err := uuid.FromString(piidStr)
	if err != nil {
		return middleware.Response{Err: errors.ErrorUnauthenticated}
	}
	data, ok := auth.Data().(*entity.AuthData)
	if !ok {
		return middleware.Response{Err: errors.ErrorUnauthenticated}
	}
	ctx := req.Context()
	for _, i := range data.Instances {
		if i.PIID == url_piid {
			ctx = context.WithValue(req.Context(), "piid", i.PIID)
			break
		}
	}
	reqWithCtx := req.WithContext(ctx)
	return next(reqWithCtx)
}
