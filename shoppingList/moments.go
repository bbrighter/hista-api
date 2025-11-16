package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
)

type MomentsParams struct {
	IfNoneMatch string `header:"If-None-Match"`
}

// encore:api auth method=GET path=/piid/:piid/moments
func (s *Service) GetMoments(ctx context.Context, piid uuid.UUID, params MomentsParams) (entity.MomentsResponse, error) {
	mom, err := s.mom.GetMoments(ctx)
	if err != nil {
		return entity.MomentsResponse{}, errors.MapError(err)
	}
	if params.IfNoneMatch == mom.ETag() {
		return mom.To304Response(), nil
	}
	list, prods, err := s.mom.GetData(ctx)
	if err != nil {
		return entity.MomentsResponse{}, errors.MapError(err)
	}
	return mom.ToResponse(list, prods), nil
}
