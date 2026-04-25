package medicines

import (
	"context"
	"slices"
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

type MedicineRepoTestSuite struct {
	suite.Suite
	ctx                context.Context
	rootDb             *gorm.DB
	tx                 *gorm.DB
	med                *MedicineRepo
	medicine1Id        uint
	medicine2Id        uint
	medicineArchivedId uint
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

var GUID = uuid.FromStringOrNil(GUID_STR)

func (suite *MedicineRepoTestSuite) SetupSuite() {
	suite.ctx = context.WithValue(context.Background(), contextKeys.Piid, GUID)
	sqlDb, err := et.NewTestDatabase(suite.ctx, "hista_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	suite.rootDb = db
	if err != nil {
		panic(err)
	}
}

func (s *MedicineRepoTestSuite) SetupTest() {
	tx := s.rootDb.Begin()
	s.tx = tx
	s.med = NewMedicineRepo(s.tx)
	s.fillWithData()
}

func (s *MedicineRepoTestSuite) TearDownSubTest() {
	err := s.tx.Exec("DELETE FROM intakes").Error
	s.Require().NoError(err)
}

func (s *MedicineRepoTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *MedicineRepoTestSuite) Debug() {
	s.tx = s.tx.Debug()
}

func (s *MedicineRepoTestSuite) fillWithData() {
	var medicine1 = Medicine{PIID: GUID, Name: "Medicine1", IsArchived: false, SortOrder: 1}
	var medicine2 = Medicine{PIID: GUID, Name: "Medicine2", IsArchived: false, SortOrder: 2}
	var archivedMedicine = Medicine{PIID: GUID, Name: "ArchivedMedicine", IsArchived: true, SortOrder: 3}

	err := s.tx.CreateInBatches([]*Medicine{&medicine1, &medicine2, &archivedMedicine}, 10).Error
	s.Require().NoError(err)
	s.medicine1Id = medicine1.ID
	s.medicine2Id = medicine2.ID
	s.medicineArchivedId = archivedMedicine.ID
}

func TestMedicineRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MedicineRepoTestSuite))
}

func (s *MedicineRepoTestSuite) TestListMedicines() {
	medicines, err := s.med.ListMedicines(s.ctx)
	s.NoError(err)
	s.Len(medicines, 3)
}

func (s *MedicineRepoTestSuite) TestCreateMedicine() {
	medicine := Medicine{Name: "Brand new", SortOrder: 20}
	err := s.med.CreateMedicine(s.ctx, &medicine)

	s.NoError(err)
	meds, err := s.med.ListMedicines(s.ctx)
	s.NoError(err)
	s.Len(meds, 4)
	slices.SortFunc(meds, func(a, b *Medicine) int {
		return a.SortOrder - b.SortOrder
	})
	med := meds[3]
	s.EqualValues(20, med.SortOrder)
	s.Equal("Brand new", med.Name)
}

func (s *MedicineRepoTestSuite) TestDeleteMedicine() {
	err := s.med.DeleteMedicine(s.ctx, s.medicine1Id)
	s.NoError(err)

	err = s.med.DeleteMedicine(s.ctx, 10000)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *MedicineRepoTestSuite) TestUpdateMedicine() {
	tests := map[string]struct {
		useWrongId    bool
		values        map[string]any
		expectError   bool
		expectedError error
	}{
		"ok":        {values: map[string]any{"name": "new name"}},
		"not found": {useWrongId: true, expectedError: gorm.ErrRecordNotFound},
		// "invalid entries": {values: map[string]any{"invalid": "a"}, expectError: true},
	}
	for name, test := range tests {
		s.Run(name, func() {
			id := s.medicine1Id
			if test.useWrongId {
				id = 1000
			}

			err := s.med.UpdateMedicine(s.ctx, id, test.values)
			if test.expectError || test.expectedError != nil {
				s.Error(err)
				if test.expectedError != nil {
					s.ErrorIs(err, test.expectedError)
				}
				return
			}
			print("assuming no error")
			s.NoError(err)
		})
	}
}

func (s *MedicineRepoTestSuite) TestBulkShiftMedicine() {
	err := s.med.BulkShiftMedicine(s.ctx)

	s.NoError(err)
	meds, err := s.med.ListMedicines(s.ctx)
	s.NoError(err)
	s.Len(meds, 3)
	slices.SortFunc(meds, func(a, b *Medicine) int { return a.SortOrder - b.SortOrder })

	s.Equal(meds[0].SortOrder, 101)
	s.Equal(meds[1].SortOrder, 102)
	s.Equal(meds[2].SortOrder, 103)
}

func (s *MedicineRepoTestSuite) TestReorderMedicines() {
	medicines := []Medicine{
		{ID: 101, PIID: GUID, Name: "Name 1", SortOrder: 100},
		{ID: 102, PIID: GUID, Name: "Name 2", SortOrder: 150},
		{ID: 103, PIID: GUID, Name: "Name 3", SortOrder: 200},
		{ID: 104, PIID: GUID, Name: "Name 4", SortOrder: 300},
		{ID: 105, PIID: GUID, Name: "Name 5", SortOrder: 250},
		{ID: 106, PIID: GUID, Name: "Name 6", SortOrder: 230},
	}
	err := s.med.db.CreateInBatches(medicines, 100).Error
	s.Require().NoError(err)

	err = s.med.ReorderMedicines(s.ctx)
	s.NoError(err)

	results, err := generic_queries.List[*Medicine](s.ctx, s.med.db)
	s.Require().NoError(err)
	slices.SortFunc(results, func(a, b *Medicine) int {
		return a.SortOrder - b.SortOrder
	})
	expectedNames := []string{
		"Medicine1", "Medicine2", "ArchivedMedicine",
		"Name 1", "Name 2", "Name 3",
		"Name 6", "Name 5", "Name 4",
	}
	expectedIds := []uint{s.medicine1Id, s.medicine2Id, s.medicineArchivedId, 101, 102, 103, 106, 105, 104}
	expectedSortOrder := []int{100, 200, 300, 400, 500, 600, 700, 800, 900}
	actualNames := []string{}
	actualIds := []uint{}
	actualSortOrder := []int{}
	for _, m := range results {
		actualNames = append(actualNames, m.Name)
		actualIds = append(actualIds, m.ID)
		actualSortOrder = append(actualSortOrder, m.SortOrder)
	}
	s.Equal(expectedNames, actualNames)
	s.Equal(expectedIds, actualIds)
	s.Equal(expectedSortOrder, actualSortOrder)
}

func (s *MedicineRepoTestSuite) TestListIntakes() {
	s.med.CreateIntake(s.ctx, &Intake{MedicineID: s.medicine1Id})

	intakes, err := s.med.ListIntakes(s.ctx)
	s.NoError(err)
	s.Len(intakes, 1)
}

func (s *MedicineRepoTestSuite) TestListGroupedIntakes() {
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
			var intakes []Intake
			for medicineId, count := range test.numberOfIntakes {
				print("Medicine ID:", medicineId)
				var date = time.Now()
				if test.useOld {
					date = time.Now().Add(-24 * 14 * time.Hour)
				}
				for range count {
					intakes = append(intakes, Intake{Date: date, MedicineID: medicineId})
				}
			}
			err := s.tx.CreateInBatches(intakes, 10).Error
			s.Require().NoError(err)

			list, err := s.med.ListGroupedIntakes(s.ctx)
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

func (s *MedicineRepoTestSuite) TestCreateIntake() {
	intake := Intake{Date: time.Now(), MedicineID: s.medicine2Id}

	err := s.med.CreateIntake(s.ctx, &intake)
	s.NoError(err)
	s.NotEqualValues(0, intake.ID)
}

func (s *MedicineRepoTestSuite) TestRemoveLastIntake() {
	tests := map[string]struct {
		medicineId    uint
		expectedError error
	}{
		"ok":        {medicineId: s.medicine1Id},
		"not found": {medicineId: 1000, expectedError: gorm.ErrRecordNotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var intake = Intake{MedicineID: s.medicine1Id}
			err := s.tx.Create(&intake).Error
			s.Require().NoError(err)

			err = s.med.RemoveLastIntake(s.ctx, test.medicineId)
			if test.expectedError != nil {
				s.ErrorIs(err, test.expectedError)
			}

		})
	}
}
