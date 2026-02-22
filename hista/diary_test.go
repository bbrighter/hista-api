package hista

import "encore.app/hista/entity"

func (s *ApiTestSuite) TestGetDiary() {
	tests := map[string]struct {
		createCondition bool
		createFood      bool
		createNote      bool
		createPollen    bool
		createIntake    bool
		expectedLen     int
	}{
		"all":       {createCondition: true, createFood: true, createNote: true, createPollen: true, createIntake: true, expectedLen: 10},
		"none":      {},
		"food":      {createFood: true, expectedLen: 1},
		"note":      {createNote: true, expectedLen: 1},
		"condition": {createCondition: true, expectedLen: 1},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.createCondition {
				s.createTestCondition()
			}
			if test.createFood {
				s.createTestFood()
			}
			if test.createNote {
				s.createTestNote()
			}
			if test.createPollen {
				s.createTestPollen()
			}
			if test.createIntake {
				s.createTestIntake()
			}
			resp, err := s.service.GetDiary(s.ctx, s.piid)
			s.NoError(err)
			s.Len(resp.Diaries, test.expectedLen)
			var types []entity.DiaryType
			for _, d := range resp.Diaries {
				types = append(types, d.Type)
			}
			if test.createCondition {
				s.Contains(types, entity.DiarySymptom)
			}
			if test.createFood {
				s.Contains(types, entity.DiaryFood)
			}
			if test.createNote {
				s.Contains(types, entity.DiaryNote)
			}
			if test.createPollen {
				s.Contains(types, entity.DiaryPollen)
			}
			if test.createIntake {
				s.Contains(types, entity.DiaryIntake)
			}
		},
		)
	}
}
