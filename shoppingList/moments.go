package shoppinglist

import (
	"context"
	"strings"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
	"encore.dev/rlog"
	"encore.dev/types/uuid"
)

type MomentsParams struct {
	IfNoneMatch string `header:"If-None-Match"`
}

// encore:api auth method=GET path=/piid/:piid/moments
func (s *Service) GetMoments(ctx context.Context, piid uuid.UUID, params MomentsParams) (entity.MomentsResponse, error) {
	rlog.Info("GET moments called...")
	mom, err := s.mom.GetMoments(ctx)
	if err != nil {
		rlog.Error("unecpted error in s.mom.GetMoments")
		return entity.MomentsResponse{}, errors.MapError(err)
	}
	rlog.Info("Got moments:", "etag", mom.ETag(), "ifnonematch", params.IfNoneMatch)
	if normalizeETag(params.IfNoneMatch) == normalizeETag(mom.ETag()) {
		return mom.To304Response(), nil
	}
	list, prods, err := s.mom.GetData(ctx)
	rlog.Info("Got the data")
	if err != nil {
		return entity.MomentsResponse{}, errors.MapError(err)
	}
	rlog.Info("momRespEtag:", "etag", mom.ToResponse(list, prods).ETag)
	return mom.ToResponse(list, prods), nil
}

func normalizeETag(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "W/")
	s = strings.Trim(s, `"`)
	return s
}
