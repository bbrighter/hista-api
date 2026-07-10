package meals

import (
	"context"
	"testing"
	"time"

	"encore.app/shared/contextKeys"
	"encore.app/shared/generic_queries"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MealRepoTestSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB // DB connection
	tx   *gorm.DB // Transaction
	m    *MealRepository
	t    *templateRepo
	i    *ingredientRepo
	piid uuid.UUID
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *MealRepoTestSuite) SetupSuite() {
	suite.piid = uuid.FromStringOrNil(GUID_STR)
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, suite.piid)
	sqlDb, err := et.NewTestDatabase(suite.ctx, "hista_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	suite.db = db
	if err != nil {
		panic(err)
	}
}

func (suite *MealRepoTestSuite) Debug() {
	suite.tx = suite.tx.Debug()
}

func (suite *MealRepoTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.Require().NoError(suite.tx.Error)
	suite.m = NewMealRepository(suite.tx)
	suite.t = newTemplateRepo(suite.tx)
	suite.i = newIngredientRepo(suite.tx)
}

func (s *MealRepoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func (s *MealRepoTestSuite) createTestMeal() uint {
	id, err := s.m.CreateMeal(s.ctx, &Meal{})
	s.Require().NoError(err)
	return id
}

func (s *MealRepoTestSuite) createTestFood(mealId uint, ingredientId uint) uint {
	food := Food{MealID: mealId, IngredientID: ingredientId}
	err := s.m.CreateFood(s.ctx, &food)
	s.Require().NoError(err)

	return food.ID
}

func (s *MealRepoTestSuite) createTestIngredient() uint {
	var ingredient = Ingredient{Name: "Name"}
	err := generic_queries.Create(s.ctx, s.tx, &ingredient)
	s.Require().NoError(err)
	return ingredient.ID
}

func (s *MealRepoTestSuite) createFullMeal() (uint, uint, uint) {
	mealId := s.createTestMeal()
	ingId := s.createTestIngredient()
	foodId := s.createTestFood(mealId, ingId)
	return mealId, foodId, ingId
}

func (s *MealRepoTestSuite) createTemplate(ingredientId uint) uint {
	var template = Template{Name: "Name", Items: []TemplateItem{
		{Condition: Raw, IngredientID: ingredientId},
	}}
	err := s.t.CreateTemplate(s.ctx, &template)
	s.Require().NoError(err)
	return template.ID
}

func TestMealRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MealRepoTestSuite))
}

func (s *MealRepoTestSuite) TestListAllMeals() {
	var meals []*Meal
	var err error

	meals, err = s.m.ListAllMeals(s.ctx)
	s.NoError(err)
	s.Len(meals, 0)

	s.createTestMeal()
	meals, err = s.m.ListAllMeals(s.ctx)
	s.NoError(err)
	s.Len(meals, 1)
	s.Len(meals[0].Foods, 0)
}

func (s *MealRepoTestSuite) TestListAllMealsAndFoodsIngredients() {
	var meals []*Meal
	var err error

	meals, err = s.m.ListMealsAndFoodsAndIngredients(s.ctx)
	s.NoError(err)
	s.Len(meals, 0)

	s.createTestFood(s.createTestMeal(), s.createTestIngredient())
	meals, err = s.m.ListMealsAndFoodsAndIngredients(s.ctx)
	s.NoError(err)
	s.Len(meals, 1)
	s.Len(meals[0].Foods, 1)
	food := meals[0].Foods[0]
	s.Equal("Name", food.Ingredient.Name)
}

func (s *MealRepoTestSuite) TestGetMealAndFoods() {
	var meal *Meal
	var err error

	_, err = s.m.GetMealAndFoods(s.ctx, 100)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	mealId, _, _ := s.createFullMeal()
	meal, err = s.m.GetMealAndFoods(s.ctx, mealId)

	s.NoError(err)
	s.Len(meal.Foods, 1)
}

func (s *MealRepoTestSuite) TestCreateMeal() {
	var meal = &Meal{}
	id, err := s.m.CreateMeal(s.ctx, meal)
	s.NoError(err)
	s.NotEqualValues(0, id)

	_, err = s.m.GetMealAndFoods(s.ctx, id)
	s.NoError(err)
}

func (s *MealRepoTestSuite) TestDeleteMeal() {
	id := s.createTestMeal()

	err := s.m.DeleteMeal(s.ctx, id)
	s.NoError(err)

	meals, err := s.m.ListAllMeals(s.ctx)
	s.NoError(err)
	s.Len(meals, 0)
}

func (s *MealRepoTestSuite) TestUpdateMeal() {
	id := s.createTestMeal()

	err := s.m.UpdateMeal(s.ctx, id, map[string]any{"stress_level": 3})
	s.NoError(err)

	meal, err := s.m.GetMealAndFoods(s.ctx, id)
	s.NoError(err)
	s.EqualValues(3, meal.StressLevel)
}

func (s *MealRepoTestSuite) TestListFoodsByMeal() {
	foods, err := s.m.ListFoodsByMeal(s.ctx, 1000)
	s.Len(foods, 0)
	s.NoError(err)

	id, _, _ := s.createFullMeal()
	newFoods, err := s.m.ListFoodsByMeal(s.ctx, id)
	s.NoError(err)
	s.Len(newFoods, 1)
}

func (s *MealRepoTestSuite) TestListFoodsByIds() {
	_, id, _ := s.createFullMeal()

	foods, err := s.m.ListFoodsByIds(s.ctx, []uint{id})
	s.NoError(err)
	s.Len(foods, 1)

	foods, err = s.m.ListFoodsByIds(s.ctx, []uint{1000})
	s.NoError(err)
	s.Len(foods, 0)
}

func (s *MealRepoTestSuite) TestCreateFoodAndIngredient() {
	// No meal
	err := s.m.CreateFoodAndIngredient(s.ctx, &Food{}, "ingredient")
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)

	// ok
	id := s.createTestMeal()
	var food = Food{MealID: id}
	err = s.m.CreateFoodAndIngredient(s.ctx, &food, "ingredient")
	s.NoError(err)
	s.NotEqualValues(0, food.ID)
	s.NotEqualValues(0, food.IngredientID)

	// Duplicate ingredient name
	err = s.m.CreateFoodAndIngredient(s.ctx, &Food{MealID: id}, "ingredient")
	s.ErrorIs(err, gorm.ErrDuplicatedKey)
}

func (s *MealRepoTestSuite) TestCreateFoodNoMeal() {
	err := s.m.CreateFood(s.ctx, &Food{})
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *MealRepoTestSuite) TestCreateFoodNoIngredient() {
	mealId := s.createTestMeal()
	err := s.m.CreateFood(s.ctx, &Food{MealID: mealId})
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}
func (s *MealRepoTestSuite) TestCreateFoodOk() {
	mealId, _, ingId := s.createFullMeal()
	food := Food{MealID: mealId, IngredientID: ingId}
	err := s.m.CreateFood(s.ctx, &food)
	s.NoError(err)

	foods, err := s.m.ListFoodsByIds(s.ctx, []uint{food.ID})
	s.NoError(err)
	s.Len(foods, 1)
	s.Equal(foods[0].IngredientID, ingId)
}

func (s *MealRepoTestSuite) TestDeleteFood() {
	var err error

	err = s.m.DeleteFood(s.ctx, 1000)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	_, foodId, _ := s.createFullMeal()
	err = s.m.DeleteFood(s.ctx, foodId)
	s.NoError(err)
	foods, err := s.m.ListFoodsByIds(s.ctx, []uint{foodId})
	s.NoError(err)
	s.Len(foods, 0)

	ings, err := s.i.ListIngredients(s.ctx)
	s.NoError(err)
	s.Len(ings, 0)

}

func (s *MealRepoTestSuite) TestDeleteFoodButOneExistsStill() {
	mealId, foodId, ingredientId := s.createFullMeal()
	s.createTestFood(mealId, ingredientId)

	err := s.m.DeleteFood(s.ctx, foodId)

	s.NoError(err)
	ings, err := s.i.ListIngredients(s.ctx)
	s.NoError(err)
	s.Len(ings, 1)
}

func (s *MealRepoTestSuite) TestFirstFood() {
	var err error

	_, err = s.m.FirstFood(s.ctx, 1000)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	_, foodId, ingId := s.createFullMeal()
	food, err := s.m.FirstFood(s.ctx, foodId)
	s.NoError(err)
	s.Equal(foodId, food.ID)
	s.Equal(ingId, food.IngredientID)
}

func (s *MealRepoTestSuite) TestUpdateFoodNotFound() {
	var err error

	err = s.m.UpdateFood(s.ctx, 1000, map[string]any{"amount": 0})
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestUpdateFoodInvalid() {
	_, foodId, _ := s.createFullMeal()
	err := s.m.UpdateFood(s.ctx, foodId, map[string]any{"a": "b"})
	s.Error(err)
}

func (s *MealRepoTestSuite) TestUpdateFood() {
	_, foodId, _ := s.createFullMeal()
	var amount int = 18
	err := s.m.UpdateFood(s.ctx, foodId, map[string]any{
		"condition": Raw,
		"amount":    &amount,
	})
	s.NoError(err)
	food, err := s.m.FirstFood(s.ctx, foodId)
	s.NoError(err)
	s.Equal(Raw, food.Condition)
	s.Equal(amount, *food.Amount)

	err = s.m.UpdateFood(s.ctx, foodId, map[string]any{
		"amount": nil,
	})
	s.NoError(err)
	food, err = s.m.FirstFood(s.ctx, foodId)
	s.Nil(food.Amount)
}

func (s *MealRepoTestSuite) TestCreateFoodsNoMealOrIng() {
	var foods = []Food{{MealID: 100, IngredientID: 10}}
	err := s.m.CreateFoods(s.ctx, foods)
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *MealRepoTestSuite) TestCreateFoods() {
	mealId, _, ingId := s.createFullMeal()
	var foods = []Food{
		{MealID: mealId, IngredientID: ingId, Condition: Raw},
		{MealID: mealId, IngredientID: ingId, Condition: Cooked},
	}

	err := s.m.CreateFoods(s.ctx, foods)
	s.NoError(err)

	foods, err = s.m.ListFoodsByMeal(s.ctx, mealId)
	s.NoError(err)
	s.Len(foods, 3)
}

func (s *MealRepoTestSuite) TestListIngredients() {
	s.createFullMeal()

	ings, err := s.i.ListIngredients(s.ctx)
	s.NoError(err)
	s.Len(ings, 1)
}

func (s *MealRepoTestSuite) TestUpdateIngredientNotFound() {
	var values = map[string]any{
		"name": "new name",
	}

	err := s.i.UpdateIngredient(s.ctx, 100, values)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}
func (s *MealRepoTestSuite) TestUpdateIngredient() {
	var values = map[string]any{
		"name": "new name",
	}

	_, _, id := s.createFullMeal()
	err := s.i.UpdateIngredient(s.ctx, id, values)
	s.NoError(err)

	ings, err := s.i.ListIngredients(s.ctx)
	s.NoError(err)
	s.Len(ings, 1)
	s.EqualValues("new name", ings[0].Name)
}

func (s *MealRepoTestSuite) TestFirstTemplateAndItems() {
	ingId := s.createTestIngredient()
	id := s.createTemplate(ingId)

	template, err := s.t.FirstTemplateAndItems(s.ctx, id)

	s.NoError(err)
	s.Equal("Name", template.Name)
	s.Len(template.Items, 1)
	s.Equal(ingId, template.Items[0].IngredientID)
	s.Equal(Raw, template.Items[0].Condition)
}

func (s *MealRepoTestSuite) TestListTemplatesAndItems() {
	ingId := s.createTestIngredient()
	id := s.createTemplate(ingId)

	templates, err := s.t.ListTemplatesAndItems(s.ctx)

	s.NoError(err)
	s.Len(templates, 1)
	template := templates[0]
	s.Equal(id, template.ID)
	s.Equal("Name", template.Name)
	s.Len(template.Items, 1)
}

func (s *MealRepoTestSuite) TestDeleteTemplate() {
	id := s.createTemplate(s.createTestIngredient())

	err := s.t.DeleteTemplate(s.ctx, id)

	s.NoError(err)
	_, err = s.t.FirstTemplateAndItems(s.ctx, id)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestCreateTemplate() {
	ingId := s.createTestIngredient()
	var template = Template{Name: "Template", Items: []TemplateItem{
		{IngredientID: ingId, Condition: Raw},
		{IngredientID: ingId, Condition: Cooked},
	}}

	err := s.t.CreateTemplate(s.ctx, &template)
	s.NoError(err)

	dbTemplate, err := s.t.FirstTemplateAndItems(s.ctx, template.ID)
	s.NoError(err)
	s.Equal("Template", dbTemplate.Name)
	s.Len(dbTemplate.Items, 2)
}

func (s *MealRepoTestSuite) TestReplaceTemplateNotFound() {
	err := s.t.ReplaceTemplate(s.ctx, 100, "name", []TemplateItem{})
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestReplaceTemplate() {
	ingId := s.createTestIngredient()
	id := s.createTemplate(ingId)

	err := s.t.ReplaceTemplate(s.ctx, id, "new name", []TemplateItem{
		{IngredientID: ingId, Condition: Raw},
		{IngredientID: ingId, Condition: Cooked},
	})

	s.NoError(err)
	template, err := s.t.FirstTemplateAndItems(s.ctx, id)
	s.NoError(err)
	s.Equal("new name", template.Name)
	s.Len(template.Items, 2)
}

func (s *MealRepoTestSuite) TestReplaceTemplateInvalidIngredient() {
	id := s.createTemplate(s.createTestIngredient())

	err := s.t.ReplaceTemplate(s.ctx, id, "new name", []TemplateItem{
		{IngredientID: 1000, Condition: Raw},
	})

	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *MealRepoTestSuite) TestAggregateNutrition() {
	var (
		proteinVal float32 = 10.0
		fiberVal   float32 = 2.0
		carbsVal   float32 = 5.2
		fatVal     float32 = 1.0
		amountVal  int     = 95
	)

	meal1Date := time.Date(2022, 3, 3, 12, 0, 0, 0, time.UTC)
	meal1 := Meal{Date: meal1Date}
	meal2 := Meal{Date: meal1Date.Add(23 * time.Hour)}
	generic_queries.Create(s.ctx, s.tx, &meal1)
	generic_queries.Create(s.ctx, s.tx, &meal2)

	ingredientWithNutrition := Ingredient{
		Name: "name",
		Nutrition: Nutrition{
			Protein:      &proteinVal,
			Fiber:        &fiberVal,
			Fat:          &fatVal,
			Carbohydrate: &carbsVal,
		},
	}
	ingredientWithoutNutrition := Ingredient{Name: "no nutrition"}
	generic_queries.Create(s.ctx, s.tx, &ingredientWithNutrition)
	generic_queries.Create(s.ctx, s.tx, &ingredientWithoutNutrition)

	foodWithAmount := Food{IngredientID: ingredientWithNutrition.ID, MealID: meal1.ID, Amount: &amountVal}
	food2 := Food{IngredientID: ingredientWithoutNutrition.ID, MealID: meal1.ID}
	food3 := Food{IngredientID: ingredientWithNutrition.ID, MealID: meal2.ID}
	gorm.G[Food](s.tx).Create(s.ctx, &foodWithAmount)
	gorm.G[Food](s.tx).Create(s.ctx, &food2)
	gorm.G[Food](s.tx).Create(s.ctx, &food3)

	tests := map[string]struct {
		fromDate      *time.Time
		toDate        *time.Time
		interval      string
		expectedError bool
		expectedLen   int
	}{
		"no dates":           {interval: "hour", expectedLen: 1},
		"invalid interval":   {interval: "something", expectedError: true},
		"from date":          {interval: "hour", fromDate: new(meal1Date.Add(-time.Minute)), expectedLen: 1},
		"from date too late": {interval: "hour", fromDate: new(meal1Date.Add(time.Minute)), expectedLen: 0},
		"to date":            {interval: "hour", toDate: new(meal1Date.Add(24 * time.Hour)), expectedLen: 1},
		"to date too early":  {interval: "hour", toDate: new(meal1Date.Add(-24 * time.Hour)), expectedLen: 0},
	}
	for name, test := range tests {
		s.Run(name, func() {
			nutrition, err := s.m.AggregateNutrition(s.ctx, test.interval, test.fromDate, test.toDate)
			if test.expectedError {
				s.Error(err)
				return
			}
			s.Len(nutrition, test.expectedLen)
			if test.expectedLen == 0 {
				return
			}

			scale := float32(amountVal) / 100.0
			expectedProtein := proteinVal * scale
			expectedFat := fatVal * scale
			expectedFiber := fiberVal * scale
			expectedCarbs := carbsVal * scale

			s.Equal(meal1Date, nutrition[0].Date.UTC())
			s.InDelta(expectedProtein, *nutrition[0].Nutrition.Protein, 0.001)
			s.InDelta(expectedFat, *nutrition[0].Nutrition.Fat, 0.001)
			s.InDelta(expectedFiber, *nutrition[0].Nutrition.Fiber, 0.001)
			s.InDelta(expectedCarbs, *nutrition[0].Nutrition.Carbohydrate, 0.001)
		})
	}
}
