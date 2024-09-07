package internal

import (
	"time"

	"encore.app/entity"
)

type ConditionEventUseCase struct {
	events IConditionEventRepo
	cats   ISymptomCategoriesRepo
}

func NewConditionEventUseCase(event IConditionEventRepo, cats ISymptomCategoriesRepo) ConditionEventUseCase {
	return ConditionEventUseCase{events: event, cats: cats}
}

func (uc ConditionEventUseCase) List() entity.ConditionEvents {
	return uc.events.ListConditionEvents()

}
func (uc ConditionEventUseCase) Create(date time.Time) (entity.ConditionEvent, error) {
	return uc.events.CreateConditionEvent(date)

}
func (uc ConditionEventUseCase) Get(id uint) (entity.ConditionEvent, error) {
	return uc.events.GetConditionEvent(id)

}
func (uc ConditionEventUseCase) Patch(id uint, date time.Time) error {
	return uc.events.PatchConditionEvent(id, date)

}
func (uc ConditionEventUseCase) Delete(id uint) (cats entity.SymptomCategories, err error) {
	err = uc.events.DeleteConditionEvent(id)
	if err == nil {
		cats = uc.cats.ListCategories()
	}
	return cats, err
}
