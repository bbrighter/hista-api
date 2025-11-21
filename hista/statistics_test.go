package hista

import (
	"time"
)

func (s *ApiTestSuite) TestGetStatisticsBySymptomIds() {
	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	ids := []uint{1}
	_, err := s.service.GetStatisticsBySymptomIds(s.ctx, s.piid, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	s.NoError(err)

	s.createTestFood()
	_, symptomId, _ := s.createTestCondition()
	ids = append(ids, symptomId)

	resp, err := s.service.GetStatisticsBySymptomIds(s.ctx, s.piid, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	s.NoError(err)
	s.Require().Len(resp.Statistics, 1)
	stat := resp.Statistics[0]
	s.EqualValues(1, stat.Count)
	s.EqualValues(1, stat.Hours1)
	s.EqualValues(1, stat.Hours24)
	s.EqualValues(1, stat.Hours72)
}

func (s *ApiTestSuite) TestGetStatisticsByIngredientsIds() {

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	ids := []uint{1}
	_, err := s.service.GetStatisticsByIngredientsIds(s.ctx, s.piid, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	s.NoError(err)

	_, ingredientId := s.createTestFood()
	s.createTestCondition()
	ids = append(ids, ingredientId)

	resp, err := s.service.GetStatisticsByIngredientsIds(s.ctx, s.piid, StatisticParams{IDs: ids, FromDate: from, ToDate: to})
	s.NoError(err)
	s.Require().Len(resp.Statistics, 1)
	stat := resp.Statistics[0]
	s.EqualValues(1, stat.Count)
	s.EqualValues(1, stat.Hours1)
	s.EqualValues(1, stat.Hours24)
	s.EqualValues(1, stat.Hours72)
}
