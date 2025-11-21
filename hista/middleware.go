package hista

import (
	"context"
	"slices"

	"encore.app/errors"
	"encore.app/shared/contextKeys"
	"encore.app/shared/entity"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
	"encore.dev/types/uuid"
)

// encore:middleware target=all
func AddPiidMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	var internalPaths []string = []string{
		"/pollen",
		"/internal/product-instance/move/:fromPiid/:toPiid",
		"/internal/product-instance/move/:toPiid"}
	if slices.Contains(internalPaths, req.Data().Path) {
		return next(req)
	}
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
			ctx = context.WithValue(req.Context(), contextKeys.Piid, i.PIID)
			break
		}
	}
	reqWithCtx := req.WithContext(ctx)
	return next(reqWithCtx)
}
