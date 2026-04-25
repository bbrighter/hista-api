package diary

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/symptoms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type listerMock struct {
	mock.Mock
}

func (m *listerMock) ListMealsAndFoodsAndIngredients(ctx context.Context) (meals.Meals, error) {
	args := m.Called(ctx)
	return args.Get(0).(meals.Meals), args.Error(1)
}
func (m *listerMock) ListConditionEventsAndDependencies(ctx context.Context) (symptoms.ConditionEvents, error) {
	args := m.Called(ctx)
	return args.Get(0).(symptoms.ConditionEvents), args.Error(1)
}
func (m *listerMock) ListSymptomCategoriesAndSymptoms(ctx context.Context) (symptoms.SymptomCategories, error) {
	args := m.Called(ctx)
	return args.Get(0).(symptoms.SymptomCategories), args.Error(1)
}
func (m *listerMock) ListNotes(ctx context.Context) ([]*notes.Note, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*notes.Note), args.Error(1)
}
func (m *listerMock) ListPollenWithSeverity(ctx context.Context, severity int) ([]pollen.PollenEvent, error) {
	args := m.Called(ctx, severity)
	return args.Get(0).([]pollen.PollenEvent), args.Error(1)
}
func (m *listerMock) ListIntakes(ctx context.Context) ([]medicines.Intake, error) {
	args := m.Called(ctx)
	return args.Get(0).([]medicines.Intake), args.Error(1)
}

var date = time.Date(2019, 3, 12, 4, 4, 0, 0, time.UTC)
var meal = &meals.Meal{
	Date: date, Freshness: meals.Fresh, StressLevel: 3, IsAlone: true,
	Foods: []meals.Food{{Condition: meals.Cooked, Ingredient: meals.Ingredient{Name: "Ingredient"}}},
}
var condEvent = &symptoms.ConditionEvent{Date: date.Add(time.Hour), Conditions: []symptoms.Condition{
	{Severity: 3, Symptom: symptoms.Symptom{Name: "Symptom", SymptomCategoryID: 3}},
}}
var cat = symptoms.SymptomCategory{ID: 3, Name: "Category"}
var note = &notes.Note{Date: date.Add(-time.Hour), Text: "Note"}
var pollenEvent = pollen.PollenEvent{CreatedAt: date.Add(time.Hour * 2), Pollens: []pollen.Pollen{
	{Type: pollen.Esche, Intensity: pollen.SmallPollen},
}}
var intake = medicines.Intake{Date: date.Add(-time.Hour * 2), Medicine: medicines.Medicine{Name: "Medicine"}}

func TestCreateDiary(t *testing.T) {
	tests := map[string]struct {
		mealReturn   meals.Meals
		eventReturn  symptoms.ConditionEvents
		catReturn    symptoms.SymptomCategories
		notesReturn  []*notes.Note
		pollenReturn []pollen.PollenEvent
		intakeReturn []medicines.Intake
		expectedLen  int
	}{
		"ok": {},
		"everything": {
			mealReturn:   meals.Meals{meal},
			eventReturn:  symptoms.ConditionEvents{condEvent},
			catReturn:    symptoms.SymptomCategories{cat},
			notesReturn:  []*notes.Note{note},
			pollenReturn: []pollen.PollenEvent{pollenEvent},
			intakeReturn: []medicines.Intake{intake},
			expectedLen:  5,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := new(listerMock)
			l.On("ListMealsAndFoodsAndIngredients", mock.Anything).Return(test.mealReturn, nil)
			l.On("ListConditionEventsAndDependencies", mock.Anything).Return(test.eventReturn, nil)
			l.On("ListSymptomCategoriesAndSymptoms", mock.Anything).Return(test.catReturn, nil)
			l.On("ListNotes", mock.Anything).Return(test.notesReturn, nil)
			l.On("ListPollenWithSeverity", mock.Anything, 1).Return(test.pollenReturn, nil)
			l.On("ListIntakes", mock.Anything).Return(test.intakeReturn, nil)
			service := &DiaryService{l: l}

			diaries, err := service.CreateDiary(t.Context())

			l.AssertExpectations(t)
			assert.NoError(t, err)
			assert.Len(t, diaries, test.expectedLen)

			for i := 1; i < len(diaries); i++ {
				assert.True(t, !diaries[i].Date.After(diaries[i-1].Date))
			}
		})
	}
}
