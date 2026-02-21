package medicine

import (
	"context"
	"testing"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
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
	intake             *IntakeRepo
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
	}))
	suite.rootDb = db
	if err != nil {
		panic(err)
	}
}

func (s *MedicineRepoTestSuite) SetupTest() {
	tx := s.rootDb.Begin()
	s.tx = tx.Debug()
	s.med = NewMedicineRepo(s.tx)
	s.intake = NewIntakeRepo(s.tx)
	s.fillWithData()
}

func (s *MedicineRepoTestSuite) TearDownSubTest() {
	err := s.tx.Debug().Exec("DELETE FROM intakes").Error
	s.Require().NoError(err)
}

func (s *MedicineRepoTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *MedicineRepoTestSuite) fillWithData() {
	var medicine1 = entity.Medicine{PIID: GUID, Name: "Medicine1", IsArchived: false}
	var medicine2 = entity.Medicine{PIID: GUID, Name: "Medicine2", IsArchived: false}
	var archivedMedicine = entity.Medicine{PIID: GUID, Name: "ArchivedMedicine", IsArchived: true}

	err := s.tx.CreateInBatches(&entity.Medicines{&medicine1, &medicine2, &archivedMedicine}, 10).Error
	s.Require().NoError(err)
	s.medicine1Id = medicine1.ID
	s.medicine2Id = medicine2.ID
	s.medicineArchivedId = archivedMedicine.ID
}

func TestMedicineRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MedicineRepoTestSuite))
}
