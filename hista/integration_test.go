package hista

import (
	"time"

	"encore.app/hista/internal/headaches"
	"encore.dev/beta/errs"
	"encore.dev/types/option"
)

func (s *ApiTestSuite) TestHeadaches() {
	var err error
	var headachesResp HeadacheListResponse
	headachesResp, err = s.service.ListHeadaches(s.ctx, s.piid)
	s.NoError(err)
	s.Len(headachesResp.Headaches, 0)

	// Create one headache
	date := time.Date(2018, 1, 2, 3, 4, 5, 0, time.Local)
	var severity uint8 = 5
	idResp, err := s.service.PostHeadache(s.ctx, s.piid, PostHeadacheParams{Date: date, Severity: severity})
	id := idResp.ID
	s.NoError(err)
	s.Greater(id, uint(0))

	// Get headache
	headachesResp, err = s.service.ListHeadaches(s.ctx, s.piid)
	s.NoError(err)
	s.Len(headachesResp.Headaches, 1)

	headache, err := s.service.GetHeadache(s.ctx, s.piid, id)
	s.NoError(err)
	s.Equal(date, headache.Date)
	s.Equal(severity, headache.Severity)
	s.Equal([]string{}, headache.Positions)

	// Patch and verify
	newDate := time.Date(2019, 1, 2, 3, 4, 5, 0, time.Local)
	err = s.service.PatchHeadache(s.ctx, s.piid, id, PatchHeadacheParams{Date: option.Some(newDate)})
	s.NoError(err)
	var newSeverity uint8 = 1
	err = s.service.PatchHeadache(s.ctx, s.piid, id, PatchHeadacheParams{Severity: option.Some(newSeverity)})
	s.NoError(err)
	newTypes := headaches.HeadacheTypes{headaches.Dull}
	err = s.service.PatchHeadache(s.ctx, s.piid, id, PatchHeadacheParams{Types: option.Some(newTypes)})
	s.NoError(err)
	newPositions := headaches.HeadachePositions{headaches.Back, headaches.Ear}
	err = s.service.PatchHeadache(s.ctx, s.piid, id, PatchHeadacheParams{Positions: option.Some(newPositions)})
	s.NoError(err)
	newSymptoms := headaches.HeadacheSymptoms{headaches.Dizziness}
	err = s.service.PatchHeadache(s.ctx, s.piid, id, PatchHeadacheParams{Symptoms: option.Some(newSymptoms)})
	s.NoError(err)

	headache, err = s.service.GetHeadache(s.ctx, s.piid, id)
	s.NoError(err)
	s.Equal(newDate, headache.Date)
	s.Equal([]string{"back", "ear"}, headache.Positions)
	s.Equal(newSeverity, headache.Severity)
	s.Equal([]string{"dizziness"}, headache.Symptoms)
	s.Equal([]string{"dull-pressing"}, headache.Types)

	// Delete and verify
	err = s.service.DeleteHeadache(s.ctx, s.piid, idResp.ID)
	s.NoError(err)
	headache, err = s.service.GetHeadache(s.ctx, s.piid, idResp.ID)
	s.Error(err)
	headachesResp, err = s.service.ListHeadaches(s.ctx, s.piid)
	s.NoError(err)
	s.Len(headachesResp.Headaches, 0)
}

func (s *ApiTestSuite) TestSymptoms() {
	symptoms, err := s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 0)

	// Create a category
	catId, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat"})
	s.NoError(err)

	// Create a symptom
	event, err := s.service.CreateConditionEvent(s.ctx, s.piid)
	s.NoError(err)
	symptomName := "symptom"
	_, err = s.service.PostCondition(s.ctx, s.piid, event.ID, ConditionRequestParams{SymptomName: option.Some(symptomName), CategoryID: option.Some(catId.ID)})
	s.NoError(err)

	symptoms, err = s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 1)
	cat := symptoms.Categories[0]
	s.Equal("cat", cat.Name)
	s.Len(cat.Symptoms, 1)
	sym := cat.Symptoms[0]
	s.Equal("symptom", sym.Name)

	// Patch category and symptom names
	err = s.service.PatchCategoryName(s.ctx, s.piid, cat.ID, PatchCategoryNameParams{Name: "new cat"})
	s.NoError(err)
	err = s.service.PatchSymptomName(s.ctx, s.piid, sym.ID, PatchSymptomNameParams{Name: "new symptom"})
	s.NoError(err)
	symptoms, err = s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 1)
	cat = symptoms.Categories[0]
	s.Equal("new cat", cat.Name)
	s.Len(cat.Symptoms, 1)
	sym = cat.Symptoms[0]
	s.Equal("new symptom", sym.Name)

	// Change category
	cat2Id, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat2"})
	s.NoError(err)
	err = s.service.PatchSymptomCategory(s.ctx, s.piid, sym.ID, PatchSymptomCategoryParams{ToCategoryID: cat2Id.ID})
	s.NoError(err)
	symptoms, err = s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 2)
	for _, c := range symptoms.Categories {
		if c.ID == cat2Id.ID {
			s.Equal("new symptom", c.Symptoms[0].Name)
		}
	}

	// Delete
	err = s.service.DeleteSymptomCategory(s.ctx, s.piid, cat.ID)
	s.NoError(err)
	err = s.service.DeleteSymptomCategory(s.ctx, s.piid, cat2Id.ID)
	s.assertErrCode(err, errs.InvalidArgument)
}

func (s *ApiTestSuite) TestMeals() {
	mealsResp, err := s.service.ListMeals(s.ctx, s.piid)
	s.NoError(err)
	s.Len(mealsResp.Meals, 0)

	postMealResp, err := s.service.PostMeal(s.ctx, s.piid, PostMealParams{Date: time.Now()}) // Why does POSt have so many params? I don't use them.
	s.NoError(err)
	mealId := postMealResp.ID

	mealsResp, err = s.service.ListMeals(s.ctx, s.piid)
	s.NoError(err)
	s.Len(mealsResp.Meals, 1)

	mealResp, err := s.service.GetMeal(s.ctx, s.piid, mealId)
	s.NoError(err)
	s.Equal(mealId, mealResp.ID)

	foodResp, err := s.service.PostFood(s.ctx, s.piid, mealId, FoodParams{IngredientName: "ing"})
	s.NoError(err)
	foodId := foodResp.Food.ID
	s.Len(foodResp.Ingredients.Ingredients, 1)
	s.NotEqualValues(0, foodResp.Food.IngredientId)
	s.Equal("cooked", foodResp.Food.Condition)

	ingResp, err := s.service.ListIngredients(s.ctx, s.piid)
	s.NoError(err)
	s.Len(ingResp.Ingredients, 1)

	err = s.service.PatchFoodCondition(s.ctx, s.piid, foodId, PatchFoodConditionParams{Condition: "raw"})
	s.NoError(err)

	var freshness uint8 = 0
	var stressLevel uint8 = 3
	var isAlone bool = false
	err = s.service.PatchMeal(s.ctx, s.piid, mealId, PatchMealParams{Freshness: &freshness, StressLevel: &stressLevel, IsAlone: &isAlone})
	s.NoError(err)
	mealResp, err = s.service.GetMeal(s.ctx, s.piid, mealId)
	s.NoError(err)
	s.Equal("raw", mealResp.Foods[0].Condition)
	s.Equal(freshness, mealResp.Freshness)
	s.Equal(stressLevel, mealResp.StressLevel)
	s.Equal(isAlone, mealResp.IsAlone)

	ings, err := s.service.DeleteFood(s.ctx, s.piid, foodId)
	s.NoError(err)
	s.Len(ings.Ingredients, 0)

	_, err = s.service.DeleteMeal(s.ctx, s.piid, mealId)
	s.NoError(err)

	_, err = s.service.GetMeal(s.ctx, s.piid, mealId)
	s.assertErrCode(err, errs.NotFound)
}

func (s *ApiTestSuite) TestManageIngredients() {
	postMealResp, err := s.service.PostMeal(s.ctx, s.piid, PostMealParams{Date: time.Now()})
	foodResp, err := s.service.PostFood(s.ctx, s.piid, postMealResp.ID, FoodParams{IngredientName: "ing"})

	ingredients, err := s.service.ListIngredients(s.ctx, s.piid)
	s.NoError(err)
	s.Len(ingredients.Ingredients, 1)
	ing := ingredients.Ingredients[0]

	err = s.service.ArchiveIngredient(s.ctx, s.piid, ing.ID)
	s.NoError(err)

	name := option.Some("new name")
	err = s.service.PatchIngredient(s.ctx, s.piid, ing.ID, PatchIngredientParams{Name: name})
	s.NoError(err)

	archived := option.Some(true)
	err = s.service.PatchIngredient(s.ctx, s.piid, ing.ID, PatchIngredientParams{Archived: archived})
	s.NoError(err)

	nutrition := PatchNutritionParams{Protein: 100, Carbohydrate: 10, Fat: 0.2, Fiber: 2}
	nutritionParams := option.Some(nutrition)
	err = s.service.PatchIngredient(s.ctx, s.piid, ing.ID, PatchIngredientParams{Nutrition: nutritionParams})
	s.NoError(err)

	ingredients, err = s.service.ListIngredients(s.ctx, s.piid)
	s.NoError(err)
	s.Len(ingredients.Ingredients, 1)
	ing = ingredients.Ingredients[0]
	s.Equal("new name", ing.Name)
	s.True(ing.IsArchived)
	s.EqualValues(100, ing.Nutrition.Protein)
	s.EqualValues(10, ing.Nutrition.Carbohydrate)
	s.EqualValues(float32(0.2), ing.Nutrition.Fat)
	s.EqualValues(2, ing.Nutrition.Fiber)

	// Deleting foods deletes ingredients as well
	_, err = s.service.DeleteFood(s.ctx, s.piid, foodResp.Food.ID)
	s.NoError(err)

	ingredients, err = s.service.ListIngredients(s.ctx, s.piid)
	s.NoError(err)
	s.Len(ingredients.Ingredients, 0)
}
