package internal

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
)

type (
	IConditionEventRepo interface {
		ListConditionEvents(ctx context.Context) ([]*entity.ConditionEvent, error)
		CreateConditionEvent(ctx context.Context, event *entity.ConditionEvent) error
		GetConditionEvent(ctx context.Context, id uint) (entity.ConditionEvent, error)
		PatchConditionEvent(ctx context.Context, id uint, date time.Time) error
		DeleteConditionEvent(ctx context.Context, id uint) error
		ListConditionEventsAndDependencies(ctx context.Context) ([]*entity.ConditionEvent, error)
	}

	IConditionEventUseCase interface {
		List(ctx context.Context) (entity.ConditionEvents, error)
		Create(ctx context.Context) (entity.ConditionEvent, error)
		Get(ctx context.Context, id uint) (entity.ConditionEvent, error)
		Patch(ctx context.Context, id uint, date time.Time) error
		Delete(ctx context.Context, id uint) (entity.SymptomCategories, error)
	}
)

type ConditionEventUseCase struct {
	events IConditionEventRepo
	cats   ISymptomCategoriesRepo
}

func NewConditionEventUseCase(event IConditionEventRepo, cats ISymptomCategoriesRepo) ConditionEventUseCase {
	return ConditionEventUseCase{events: event, cats: cats}
}

func (uc ConditionEventUseCase) List(ctx context.Context) (entity.ConditionEvents, error) {
	return uc.events.ListConditionEvents(ctx)

}
func (uc ConditionEventUseCase) Create(ctx context.Context) (entity.ConditionEvent, error) {
	var event = &entity.ConditionEvent{Date: time.Now()}
	err := uc.events.CreateConditionEvent(ctx, event)
	return *event, errors.MapError(err)

}
func (uc ConditionEventUseCase) Get(ctx context.Context, id uint) (entity.ConditionEvent, error) {
	event, err := uc.events.GetConditionEvent(ctx, id)
	return event, errors.MapError(err)

}
func (uc ConditionEventUseCase) Patch(ctx context.Context, id uint, date time.Time) error {
	return errors.MapError(uc.events.PatchConditionEvent(ctx, id, date))

}
func (uc ConditionEventUseCase) Delete(ctx context.Context, id uint) (cats entity.SymptomCategories, err error) {
	err = uc.events.DeleteConditionEvent(ctx, id)
	if err == nil {
		cats, err = uc.cats.ListCategories(ctx)
		return cats, errors.MapError(err)
	}
	return cats, errors.MapError(err)
}
