package internal

import (
	"context"

	"encore.app/shoppingList/entity"
)

type IMomentUseCase interface {
	GetMoments(ctx context.Context) (*entity.Moment, error)
	GetData(ctx context.Context) (*entity.List, []*entity.Product, error)
	UpdateMoments(ctx context.Context) error
}

type MomentsUseCase struct {
	m MomentRepo
	l ListRepo
	p ProductRepo
}

func NewMomentsUseCase(m MomentRepo, l ListRepo, p ProductRepo) MomentsUseCase {
	return MomentsUseCase{m: m, l: l, p: p}
}

func (m MomentsUseCase) GetData(ctx context.Context) (list *entity.List, prods []*entity.Product, err error) {
	list, err = m.l.First(ctx)
	if err != nil {
		return list, prods, err
	}
	prods, err = m.p.List(ctx)

	return list, prods, err
}

func (m MomentsUseCase) GetMoments(ctx context.Context) (moment *entity.Moment, err error) {
	return m.m.GetOrCreate(ctx)
}

func (m MomentsUseCase) UpdateMoments(ctx context.Context) error {
	return m.m.Update(ctx)
}
