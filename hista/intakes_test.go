package hista

import (
	"time"
)

func (s *ApiTestSuite) TestIntakes() {
	listResp, err := s.service.ListIntakes(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Intakes, 0)

	idResp, err := s.service.CreateMedicine(s.ctx, s.piid, PostMedicineParams{Name: "Name"})
	s.NoError(err)
	medicineId := idResp.ID

	err = s.service.IncrementIntake(s.ctx, s.piid, medicineId)
	s.NoError(err)

	listResp, err = s.service.ListIntakes(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Intakes, 1)
	s.EqualValues(1, listResp.Intakes[0].Count)
	s.EqualValues(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC), listResp.Intakes[0].Date)

	err = s.service.DecrementIntake(s.ctx, s.piid, medicineId)
	s.NoError(err)

	listResp, err = s.service.ListIntakes(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Intakes, 0)
}
