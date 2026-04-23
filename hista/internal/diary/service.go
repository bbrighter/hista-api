package diary

import (
	"context"
	"sort"

	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/symptoms"
)

type Lister interface {
	ListMealsAndFoodsAndIngredients(context.Context) (meals.Meals, error)
	ListConditionEventsAndDependencies(context.Context) (symptoms.ConditionEvents, error)
	ListSymptomCategoriesAndSymptoms(context.Context) (symptoms.SymptomCategories, error)
	ListNotes(context.Context) ([]*notes.Note, error)
	ListPollenWithSeverity(context.Context, int) ([]pollen.PollenEvent, error)
	ListIntakes(context.Context) ([]medicines.Intake, error)
}

type DiaryService struct {
	l Lister
}

type listerImpl struct {
	mealsRepo     *meals.MealRepository
	symptomsRepo  *symptoms.SymptomRepo
	conditionRepo *symptoms.ConditionRepo
	notesRepo     *notes.NotesRepo
	pollenRepo    *pollen.PollenRepo
	medRepo       *medicines.MedicineRepo
}

func (l *listerImpl) ListConditionEventsAndDependencies(ctx context.Context) (symptoms.ConditionEvents, error) {
	return l.conditionRepo.ListConditionEventsAndDependencies(ctx)
}
func (l *listerImpl) ListSymptomCategoriesAndSymptoms(ctx context.Context) (symptoms.SymptomCategories, error) {
	return l.symptomsRepo.ListSymptomCategoriesAndSymptoms(ctx)
}
func (l *listerImpl) ListNotes(ctx context.Context) ([]*notes.Note, error) {
	return l.notesRepo.ListNotes(ctx)
}
func (l *listerImpl) ListIntakes(ctx context.Context) ([]medicines.Intake, error) {
	return l.medRepo.ListIntakes(ctx)
}

func (l *listerImpl) ListMealsAndFoodsAndIngredients(ctx context.Context) (meals.Meals, error) {
	return l.mealsRepo.ListMealsAndFoodsAndIngredients(ctx)
}

func (l *listerImpl) ListPollenWithSeverity(ctx context.Context, severity int) ([]pollen.PollenEvent, error) {
	return l.pollenRepo.ListPollenWithSeverity(ctx, severity)
}

func NewDiaryService(
	mealsRepo *meals.MealRepository,
	conditionRepo *symptoms.ConditionRepo,
	symptomsRepo *symptoms.SymptomRepo,
	notesRepo *notes.NotesRepo,
	pollenRepo *pollen.PollenRepo,
	medRepo *medicines.MedicineRepo,
) *DiaryService {
	return &DiaryService{
		l: &listerImpl{
			mealsRepo:     mealsRepo,
			symptomsRepo:  symptomsRepo,
			conditionRepo: conditionRepo,
			notesRepo:     notesRepo,
			pollenRepo:    pollenRepo,
			medRepo:       medRepo,
		},
	}
}
func (s *DiaryService) CreateDiary(ctx context.Context) ([]Diary, error) {
	parts := [][]Diary{}

	meals, err := s.l.ListMealsAndFoodsAndIngredients(ctx)
	if err != nil {
		return []Diary{}, err
	}
	parts = append(parts, mealsToDiary(meals))

	events, err := s.l.ListConditionEventsAndDependencies(ctx)
	if err != nil {
		return []Diary{}, err
	}
	cats, err := s.l.ListSymptomCategoriesAndSymptoms(ctx)
	if err != nil {
		return []Diary{}, err
	}
	parts = append(parts, symptomsToDiary(events, cats))

	notes, err := s.l.ListNotes(ctx)
	if err != nil {
		return []Diary{}, err
	}
	parts = append(parts, notesToDiary(notes))

	pollens, err := s.l.ListPollenWithSeverity(ctx, 1)
	if err != nil {
		return []Diary{}, err
	}
	parts = append(parts, pollenToDiary(pollens))

	intakes, err := s.l.ListIntakes(ctx)
	if err != nil {
		return []Diary{}, err
	}
	parts = append(parts, intakesToDiary(intakes))

	totalLen := 0
	for _, p := range parts {
		totalLen += len(p)
	}

	diaries := make([]Diary, 0, totalLen)
	for _, p := range parts {
		diaries = append(diaries, p...)
	}

	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})

	return diaries, nil
}
