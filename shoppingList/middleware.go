package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.app/shared/contextKeys"
	"encore.app/shared/entity"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
	"encore.dev/types/uuid"
)

// encore:middleware target=all
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
	for _, inst := range data.Instances {
		if inst.PIID == url_piid {
			ctx = context.WithValue(req.Context(), contextKeys.Piid, inst.PIID)
			if !validateApp(inst.AppMapping) {
				return middleware.Response{Err: errors.ErrorUnauthenticated}
			}
			break
		}
	}

	reqWithCtx := req.WithContext(ctx)
	return next(reqWithCtx)
}

func validateApp(appMapping map[string]bool) bool {
	return appMapping["shopping-list"]

}
