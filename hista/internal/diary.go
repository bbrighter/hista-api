package internal

import (
	"context"

	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/symptoms"
)

type IDiaryUseCase interface {
	Get(ctx context.Context) (meals.Meals, symptoms.ConditionEvents, symptoms.SymptomCategories, []*notes.Note, []pollen.PollenEvent, []medicines.Intake)
}

type DiaryUseCase struct {
	meals   *meals.MealRepository
	events  *symptoms.ConditionRepo
	cats    *symptoms.SymptomRepo
	notes   *notes.NotesRepo
	pollens *pollen.PollenRepo
	intakes *medicines.MedicineRepo
}

func NewDiaryUseCase(
	meals *meals.MealRepository,
	events *symptoms.ConditionRepo,
	cats *symptoms.SymptomRepo,
	notes *notes.NotesRepo,
	pollens *pollen.PollenRepo,
	intakes *medicines.MedicineRepo,
) DiaryUseCase {
	return DiaryUseCase{
		meals: meals, events: events, cats: cats, notes: notes, pollens: pollens, intakes: intakes}
}

func (uc DiaryUseCase) Get(ctx context.Context) (
	meals meals.Meals,
	events symptoms.ConditionEvents,
	cats symptoms.SymptomCategories,
	notes []*notes.Note,
	pollens []pollen.PollenEvent,
	intakes []medicines.Intake) {
	meals, _ = uc.meals.ListMealsAndFoodsAndIngredients(ctx)
	events, _ = uc.events.ListConditionEventsAndDependencies(ctx)
	cats, _ = uc.cats.ListSymptomCategoriesAndSymptoms(ctx)
	notes, _ = uc.notes.ListNotes(ctx)
	pollens, _ = uc.pollens.ListPollenWithSeverity(ctx, 1)
	intakes, _ = uc.intakes.ListIntakes(ctx)

	return meals, events, cats, notes, pollens, intakes
}
