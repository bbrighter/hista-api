package meals

import (
	"encore.app/hista/entity"
	"gorm.io/gorm"
)

func (s *MealRepoTestSuite) TestGetIngredients() {
	s.createIngredient()
	ing, err := s.repo.ListIngredients(s.ctx)
	s.NoError(err)
	s.Len(ing, 1)
}

func (s *MealRepoTestSuite) TestUpdateIngredientNameOk() {
	ingId := s.createIngredient()

	values := map[string]any{"name": "new name"}
	err := s.repo.UpdateIngredient(s.ctx, ingId, values)
	s.NoError(err)
}

func (s *MealRepoTestSuite) TestUpdateIngredientNameDuplicate() {
	ingId := s.createIngredient()
	s.createIngredientWithProps("existing", nil)

	values := map[string]any{"name": "existing"}
	err := s.repo.UpdateIngredient(s.ctx, ingId, values)
	s.Error(err)
	s.ErrorIs(err, gorm.ErrDuplicatedKey)
}

func (s *MealRepoTestSuite) TestUpdateNotFound() {
	err := s.repo.UpdateIngredient(s.ctx, 100, make(map[string]any))
	s.Error(err)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestUpdateNutrition() {
	ingId := s.createIngredient()

	values := map[string]any{"nutrition_protein": 100, "nutrition_fat": 30}
	err := s.repo.UpdateIngredient(s.ctx, ingId, values)
	s.NoError(err)
}

func (s *MealRepoTestSuite) TestToggleArchivedOk() {
	ingId := s.createIngredient()
	getArchivedValue := func() bool {
		ing, err := gorm.G[entity.Ingredient](s.tx).Where("id = ?", ingId).First(s.ctx)
		s.Require().NoError(err)
		return ing.IsArchived
	}

	err := s.repo.ToggleArchived(s.ctx, ingId)
	s.NoError(err)
	s.True(getArchivedValue())

	err = s.repo.ToggleArchived(s.ctx, ingId)
	s.NoError(err)
	s.False(getArchivedValue())

}

func (s *MealRepoTestSuite) TestToggleArchivedNotFound() {
	err := s.repo.ToggleArchived(s.ctx, 100)
	s.Error(err)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestDeleteIngredientStillInUse() {
	ingId := s.createIngredient()
	s.createFood(ingId)

	err := s.repo.DeleteIngredient(s.ctx, ingId)
	s.Error(err)
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *MealRepoTestSuite) TestDeleteIngredientNotFound() {
	err := s.repo.DeleteIngredient(s.ctx, 100)
	s.Error(err)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MealRepoTestSuite) TestDeleteIngredientOk() {
	ingId := s.createIngredient()
	err := s.repo.DeleteIngredient(s.ctx, ingId)
	s.NoError(err)
}

// func (s *MealRepoTestSuite) TestUpdateIngredients() {

// }
