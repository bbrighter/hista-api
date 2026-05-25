package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

type MomentsResponse struct {
	Items    int `json:"itemsVersion"`
	Products int `json:"productsVersion"`
}

// encore:api auth method=GET path=/piid/:piid/moments
func (s *Service) GetMoments(ctx context.Context, piid uuid.UUID) (MomentsResponse, error) {
	mom, err := s.sm.GetMoment(ctx)
	if err != nil {
		return MomentsResponse{}, errors.MapError(err)
	}

	return MomentsResponse{
		Items:    mom.ItemsVersion,
		Products: mom.ProductsVersion,
	}, nil
}
