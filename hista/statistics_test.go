package hista

import (
	"time"
)

func (s *ApiTestSuite) TestGetStatisticsByIngredientId() {

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	var id uint = 1
	_, err := s.service.GetStatisticsByIngredientId(s.ctx, s.piid, MealStatisticsParams{ID: id, FromDate: from, ToDate: to})
	s.NoError(err)

	_, ingredientId := s.createTestFood()
	s.createTestCondition()

	resp, err := s.service.GetStatisticsByIngredientId(s.ctx, s.piid, MealStatisticsParams{ID: ingredientId, FromDate: from, ToDate: to})
	s.NoError(err)
	s.Require().Len(resp.Statistics, 1)
	s.EqualValues(1, resp.Count)
	stat := resp.Statistics[0]
	s.EqualValues(1, stat.Hours1)
	s.EqualValues(0, stat.Hours24)
	s.EqualValues(0, stat.Hours72)
}
