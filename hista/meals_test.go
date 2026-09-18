package hista

import (
	"time"

	"encore.app/hista/internal/meals"
	"encore.dev/beta/errs"
	"gorm.io/gorm"
)

func (s *ApiTestSuite) TestGetMealsAPI() {
	resp, err := s.service.ListMeals(s.ctx, s.piid)
	s.NoError(err)
	s.Len(resp.Meals, 0)

	s.createTestMeal()

	resp, err = s.service.ListMeals(s.ctx, s.piid)
	s.NoError(err)
	s.Len(resp.Meals, 1)
}

func (s *ApiTestSuite) TestPostMealAPI() {
	var params = PostMealParams{Date: time.Now()}
	resp, err := s.service.PostMeal(s.ctx, s.piid, params)
	defer s.service.DeleteMeal(s.ctx, s.piid, resp.ID)

	s.NoError(err)
	s.GreaterOrEqual(resp.ID, uint(1))
}

func (s *ApiTestSuite) TestGetMealAPI() {
	tests := map[string]struct {
		useWrongId      bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			id := s.createTestMeal()
			if test.useWrongId {
				id = 1000
			}
			_, err := s.service.GetMeal(s.ctx, s.piid, id)
			s.assertErrCode(err, test.expectedErrCode)
		})
	}
}

func (s *ApiTestSuite) TestDeleteMealAPI() {
	tests := map[string]struct {
		useWrongId      bool
		expectedErrCode errs.ErrCode
	}{
		"ok":        {},
		"not found": {useWrongId: true, expectedErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			mealId := s.createTestMeal()
			if test.useWrongId {
				mealId = 1000
			}
			ings, err := s.service.DeleteMeal(s.ctx, s.piid, mealId)
			s.assertErrCode(err, test.expectedErrCode)
			if test.expectedErrCode == 0 {
				s.Len(ings.Ingredients, 0)
			}
		})
	}
}

func (s *ApiTestSuite) TestPatchMealAPI() {
	var params PatchMealParams
	var now time.Time = time.Now()
	params.Date = &now

	err := s.service.PatchMeal(s.ctx, s.piid, 10000, params)
	s.Error(err)
	id := s.createTestMeal()

	err = s.service.PatchMeal(s.ctx, s.piid, id, params)
	s.NoError(err)

	var stressLevel uint8 = 2
	params.StressLevel = &stressLevel
	err = s.service.PatchMeal(s.ctx, s.piid, id, params)
	s.NoError(err)

}

func (s *ApiTestSuite) TestGetFoods() {
	resp, err := s.service.GetFoods(s.ctx, s.piid, 1)
	s.NoError(err)

	id := s.createTestMeal()

	resp, err = s.service.GetFoods(s.ctx, s.piid, id)
	s.NoError(err)
	s.Len(resp.Foods, 0)
}

func (s *ApiTestSuite) TestPostFoodAPI() {

	mealId := s.createTestMeal()

	var params = FoodParams{IngredientName: "New"}
	_, err := s.service.PostFood(s.ctx, s.piid, mealId, params)
	s.NoError(err)
}

// func (s *ApiTestSuite) TestDeleteFoodAPI() {
// 	_, err := s.service.DeleteFood(s.ctx, s.piid, 100)
// 	s.EqualError(err, "not_found: not found")

// 	mealId := s.createTestMeal()
// 	var params = FoodParams{IngredientName: "New"}
// 	food, _ := s.service.PostFood(s.ctx, s.piid, mealId, params)

// 	ing, err := s.service.DeleteFood(s.ctx, s.piid, food.Food.ID)
// 	s.NoError(err)
// 	s.Len(ing.Ingredients, 0)
// }

func (s *ApiTestSuite) TestPostFoodByTemplate() {
	_, ingId := s.createTestFood()

	template, err := s.service.PostTemplate(
		s.ctx,
		s.piid,
		TemplateParams{
			Name: "Template",
			Items: []TemplateItemParams{
				{IngredientId: ingId, Condition: "cooked"},
			}},
	)
	s.NoError(err)
	s.NotEqualValues(0, template.ID)
}

func (s *ApiTestSuite) TestPostFoodByNameIsNotArchived() {
	mealId := s.createTestMeal()

	resp, err := s.service.PostFood(s.ctx, s.piid, mealId, FoodParams{IngredientName: "New name", IngredientID: 0})
	s.NoError(err)
	s.Len(resp.Ingredients.Ingredients, 1)
	s.False(resp.Ingredients.Ingredients[0].IsArchived)
}

func (s *ApiTestSuite) TestPostArchivedFoodByIdIsNotArchived() {
	mealId := s.createTestMeal()
	var ingredient = meals.Ingredient{PIID: s.piid, Name: "Name", IsArchived: true}
	err := gorm.G[meals.Ingredient](s.db).Create(s.ctx, &ingredient)
	s.Require().NoError(err)

	resp, err := s.service.PostFood(s.ctx, s.piid, mealId, FoodParams{IngredientID: ingredient.ID})
	s.NoError(err)
	s.Len(resp.Ingredients.Ingredients, 1)
	s.False(resp.Ingredients.Ingredients[0].IsArchived)
}

func (s *ApiTestSuite) TestPostNonArchivedFoodByIdIsNotArchived() {
	mealId := s.createTestMeal()
	var ingredient = meals.Ingredient{PIID: s.piid, Name: "Name", IsArchived: false}
	err := gorm.G[meals.Ingredient](s.db).Create(s.ctx, &ingredient)
	s.Require().NoError(err)

	resp, err := s.service.PostFood(s.ctx, s.piid, mealId, FoodParams{IngredientID: ingredient.ID})
	s.NoError(err)
	s.Len(resp.Ingredients.Ingredients, 1)
	s.False(resp.Ingredients.Ingredients[0].IsArchived)
}
