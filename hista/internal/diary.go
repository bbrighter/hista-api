package internal

import (
	"context"

	"encore.app/hista/entity"
)

type IDiaryUseCase interface {
	Get(ctx context.Context) (entity.Meals, entity.ConditionEvents, entity.SymptomCategories, entity.Notes, entity.PollenEvents, []*entity.Intake)
}

type DiaryUseCase struct {
	meals   IMealsRepository
	events  IConditionEventRepo
	cats    ISymptomCategoriesRepo
	notes   INotesRepository
	pollens IPollenRepo
	intakes IIntakeRepo
}

func NewDiaryUseCase(meals IMealsRepository,
	events IConditionEventRepo,
	cats ISymptomCategoriesRepo,
	notes INotesRepository,
	pollens IPollenRepo,
	intakes IIntakeRepo,
) DiaryUseCase {
	return DiaryUseCase{
		meals: meals, events: events, cats: cats, notes: notes, pollens: pollens, intakes: intakes}
}

func (uc DiaryUseCase) Get(ctx context.Context) (
	meals entity.Meals,
	events entity.ConditionEvents,
	cats entity.SymptomCategories,
	notes entity.Notes,
	pollens entity.PollenEvents,
	intakes []*entity.Intake) {
	meals, _ = uc.meals.ListMealsWithDependencies(ctx)
	events, _ = uc.events.ListConditionEventsAndDependencies(ctx)
	cats, _ = uc.cats.ListCategories(ctx)
	notes, _ = uc.notes.List(ctx)
	pollens = uc.pollens.FindPollenWithSeverity(1)
	intakes, _ = uc.intakes.List(ctx)

	return meals, events, cats, notes, pollens, intakes
}
