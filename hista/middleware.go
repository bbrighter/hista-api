package hista

import (
	"context"
	"strings"
	"time"

	"encore.app/errors"
	"encore.app/shared/contextKeys"
	"encore.app/shared/entity"
	"encore.dev/beta/auth"
	"encore.dev/metrics"
	"encore.dev/middleware"
	"encore.dev/types/uuid"
)

// encore:middleware target=all
func AddPiidMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	if req.Data().Path == "/pollen" {
		return next(req)
	}
	if strings.HasPrefix(req.Data().Path, "/internal/product-instance/move/") {
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

var ResponseTime = metrics.NewGauge[int64]("response_time_ms", metrics.GaugeConfig{})

// encore:middleware target=all
func MeasureExecutionTime(req middleware.Request, next middleware.Next) middleware.Response {
	start := time.Now()

	resp := next(req)

	duration := time.Since(start)

	ResponseTime.Set(duration.Milliseconds())
	return resp
}
