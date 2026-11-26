package hista

import (
	"context"

	"encore.app/hista/internal/meta"
	"encore.dev/types/uuid"
)

type MetaResponse struct {
	EntitiesFilled []EntitiesFilledResponse `json:"entitiesFilled"`
}

type EntitiesFilledResponse struct {
	Entity string `json:"entity"`
	Filled bool   `json:"filled"`
}

func entitiesToMetaResponse(entities []meta.CountEntities) MetaResponse {
	var entitiesFilled = []EntitiesFilledResponse{}
	for _, entity := range entities {
		entitiesFilled = append(entitiesFilled,
			EntitiesFilledResponse{
				Entity: entity.Entity,
				Filled: entity.Filled,
			})
	}
	return MetaResponse{EntitiesFilled: entitiesFilled}
}

// encore:api auth method=GET path=/piid/:piid/meta
func (s *Service) GetMetaInformation(ctx context.Context, piid uuid.UUID) (MetaResponse, error) {
	entities, err := s.meta.MetaCountInformation(ctx)
	if err != nil {
		return MetaResponse{}, err
	}

	return entitiesToMetaResponse(entities), nil
}
