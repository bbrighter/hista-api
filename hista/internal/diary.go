package internal

import (
	"context"

	"encore.app/hista/entity"
)

type IDiaryUseCase interface {
	Get(ctx context.Context) (entity.Meals, entity.ConditionEvents, entity.SymptomCategories, entity.Notes, entity.PollenEvents)
}

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

func (uc DiaryUseCase) Get(ctx context.Context) (
	meals entity.Meals,
	events entity.ConditionEvents,
	cats entity.SymptomCategories,
	notes entity.Notes,
	pollens entity.PollenEvents) {
	meals, _ = uc.meals.ListMealsWithDependencies(ctx)
	events, _ = uc.events.ListConditionEventsAndDependencies(ctx)
	cats, _ = uc.cats.ListCategories(ctx)
	notes, _ = uc.notes.List(ctx)
	pollens = uc.pollens.FindPollenWithSeverity(1)

	return meals, events, cats, notes, pollens
}
