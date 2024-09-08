package internal

import (
	"encore.app/entity"
)

type DiaryUseCase struct {
	meals   IMealsRepository
	events  IConditionEventRepo
	cats    ISymptomCategoriesRepo
	notes   INotesRepository
	pollens IPollenRepo
}

func NewDiaryUseCase(meals IMealsRepository,
	events IConditionEventRepo,
	cats ISymptomCategoriesRepo,
	notes INotesRepository,
	pollens IPollenRepo) DiaryUseCase {
	return DiaryUseCase{
		meals: meals, events: events, cats: cats, notes: notes, pollens: pollens}
}

func (uc DiaryUseCase) Get() (
	meals entity.Meals,
	events entity.ConditionEvents,
	cats entity.SymptomCategories,
	notes entity.Notes,
	pollens entity.PollenEvents) {
	meals = uc.meals.ListMealsWithDependencies()
	events = uc.events.ListConditionEventsAndDependencies()
	cats = uc.cats.ListCategories()
	notes = uc.notes.List()
	pollens = uc.pollens.FindPollenWithSeverity(1)

	return meals, events, cats, notes, pollens
}
