package meta

import "context"

type MetaService struct {
	m *MetaRepo
}

func NewMetaService(m *MetaRepo) *MetaService {
	return &MetaService{m: m}
}

func (s *MetaService) MetaCountInformation(ctx context.Context) ([]CountEntities, error) {
	return s.m.AreTodaysEntitiesFilled(ctx)
}
