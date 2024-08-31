package internal

import (
	"time"

	"encore.app/entity"
)

type ConditionEventUseCase struct {
	events IConditionEventRepo
}

func NewConditionEventUseCase(event IConditionEventRepo) ConditionEventUseCase {
	return ConditionEventUseCase{events: event}
}

func (uc ConditionEventUseCase) List() entity.ConditionEvents {
	return uc.events.ListConditionEvents()

}
func (uc ConditionEventUseCase) Create(date time.Time) (uint, error) {
	return uc.events.CreateConditionEvent(date)

}
func (uc ConditionEventUseCase) Get(id uint) (entity.ConditionEvent, error) {
	return uc.events.GetConditionEvent(id)

}
func (uc ConditionEventUseCase) Patch(id uint, date time.Time) error {
	return uc.events.PatchConditionEvent(id, date)

}
func (uc ConditionEventUseCase) Delete(id uint) error {
	return uc.events.DeleteConditionEvent(id)

}
