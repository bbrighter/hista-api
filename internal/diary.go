package internal

import (
	"encore.app/entity"
	"encore.app/pollen"
)

type DiaryUseCase struct {
	meals  IMealsRepository
	events IConditionEventRepo
	cats   ISymptomCategoriesRepo
	notes  INotesRepository
}

func NewDiaryUseCase(meals IMealsRepository,
	events IConditionEventRepo,
	cats ISymptomCategoriesRepo,
	notes INotesRepository) DiaryUseCase {
	return DiaryUseCase{
		meals: meals, events: events, cats: cats, notes: notes}
}

func (uc DiaryUseCase) Get() (
	meals entity.Meals,
	events entity.ConditionEvents,
	cats entity.SymptomCategories,
	notes entity.Notes,
	pollens pollen.PollenEvents) {
	meals = uc.meals.ListMealsWithDependencies()
	events = uc.events.ListConditionEventsAndDependencies()
	cats = uc.cats.ListCategories()
	notes = uc.notes.List()
	// TODO POLLENS

	return meals, events, cats, notes, pollens
}
