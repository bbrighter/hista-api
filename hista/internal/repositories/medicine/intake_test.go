package medicine

import (
	"time"

	"encore.app/hista/entity"
	"gorm.io/gorm"
)

func (s *MedicineRepoTestSuite) TestListGrouped() {
	tests := map[string]struct {
		numberOfIntakes   map[uint]int
		useOld            bool
		expectedLength    int
		countForMedicine1 int64
	}{
		"empty list":      {},
		"1 intake, med1":  {numberOfIntakes: map[uint]int{s.medicine1Id: 1}, expectedLength: 1, countForMedicine1: 1},
		"2 intakes, med1": {numberOfIntakes: map[uint]int{s.medicine1Id: 2}, expectedLength: 1, countForMedicine1: 2},
		"old not shown":   {numberOfIntakes: map[uint]int{s.medicine1Id: 1}, expectedLength: 0, countForMedicine1: 0, useOld: true},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var intakes []entity.Intake
			for medicineId, count := range test.numberOfIntakes {
				var date = time.Now()
				if test.useOld {
					date = time.Now().Add(-24 * 14 * time.Hour)
				}
				for range count {
					intakes = append(intakes, entity.Intake{PIID: GUID, Date: date, MedicineID: medicineId, MedicinePIID: GUID})
				}
			}
			err := s.tx.CreateInBatches(intakes, 10).Error
			s.Require().NoError(err)

			list, err := s.intake.ListGrouped(s.ctx)
			s.NoError(err)
			s.Len(list, test.expectedLength)

			for _, l := range list {
				if l.MedicineId == s.medicine1Id {
					s.Equal(test.countForMedicine1, l.Count)
				}
			}

		})
	}
}

func (s *MedicineRepoTestSuite) TestRemove() {
	tests := map[string]struct {
		medicineId    uint
		expectedError error
	}{
		"ok":        {medicineId: s.medicine1Id},
		"not found": {medicineId: 1000, expectedError: gorm.ErrRecordNotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var intake = entity.Intake{PIID: GUID, MedicineID: s.medicine1Id, MedicinePIID: GUID}
			err := s.tx.Create(&intake).Error
			s.Require().NoError(err)

			err = s.intake.Remove(s.ctx, test.medicineId)
			if test.expectedError != nil {
				s.ErrorIs(err, test.expectedError)
			}

		})
	}
}
