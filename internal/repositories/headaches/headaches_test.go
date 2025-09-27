package headaches

import (
	"context"
	"testing"
	"time"

	"encore.app/entity"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const TEST_GUID = "0c5e945e-ef6c-4934-91ff-702d94e2e7a8"

func initTest(t *testing.T, addHeadache bool) (*HeadacheRepository, context.Context, entity.Headache) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	var err error
	err = db.AutoMigrate(&entity.Headache{})
	assert.NoError(t, err)

	var headache entity.Headache
	if addHeadache {
		headache = entity.Headache{
			ID:        1,
			Date:      time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			Severity:  3,
			Types:     []entity.HeadacheType{entity.Dull},
			Positions: []entity.HeadachePosition{entity.Front, entity.Back},
			Symptoms:  []entity.HeadacheSymptom{entity.ConcentrationLack},
			PIID:      uuid.FromStringOrNil(TEST_GUID),
		}
		err = db.Create(&headache).Error
		assert.NoError(t, err)
	}

	ctx := t.Context()
	ctx = context.WithValue(ctx, "piid", uuid.FromStringOrNil(TEST_GUID))

	return NewHeadacheRepository(db), ctx, headache
}
func TestListHeadaches(t *testing.T) {
	tests := map[string]struct {
		addHeadache      bool
		expectedLen      int
		expectedSeverity int
	}{
		"no headaches": {
			addHeadache: false,
			expectedLen: 0,
		},
		"one headache": {
			addHeadache:      true,
			expectedLen:      1,
			expectedSeverity: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx, _ := initTest(t, test.addHeadache)

			headaches, err := repo.ListHeadaches(ctx)
			assert.NoError(t, err)
			require.Len(t, headaches, test.expectedLen)
			if test.expectedLen > 0 {
				assert.EqualValues(t, test.expectedSeverity, headaches[0].Severity)
			}
		})
	}
}

func TestCreateHeadache(t *testing.T) {
	tests := map[string]struct {
		headache      *entity.Headache
		errorExpected bool
	}{
		"valid headache": {
			headache:      &entity.Headache{Date: time.Now(), Severity: 3, Positions: []entity.HeadachePosition{entity.Front}},
			errorExpected: false,
		},
		"invalid headache": {
			headache:      &entity.Headache{Positions: []entity.HeadachePosition{"something"}},
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx, _ := initTest(t, false)

			err := repo.CreateHeadache(ctx, test.headache)
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Greater(t, test.headache.ID, uint(0), "id > 0 is returned")
			}
		})
	}
}

func TestDeleteHeadache(t *testing.T) {
	tests := map[string]struct {
		haId    uint
		isError bool
	}{
		"found":     {haId: 1},
		"not found": {haId: 10, isError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx, _ := initTest(t, true)
			err := repo.DeleteHeadache(ctx, test.haId)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetHeadache(t *testing.T) {
	tests := map[string]struct {
		haId  uint
		found bool
	}{
		"found":     {haId: 1, found: true},
		"not found": {haId: 10, found: false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx, origH := initTest(t, true)
			headache, err := repo.GetHeadache(ctx, test.haId)
			if test.found {
				assert.NoError(t, err)
				assert.Equal(t, &origH, headache)
			} else {
				assert.Error(t, err)
				return
			}
		})
	}
}

func TestPatchHeadache(t *testing.T) {
	date := time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)
	var severity entity.HeadacheSeverity = 2
	types := entity.HeadacheTypes{entity.Pulsating}
	positions := entity.HeadachePositions{entity.Temple}
	symptoms := entity.HeadacheSymptoms{entity.Tinnitus}
	description := "desc"
	tests := map[string]struct {
		haID        uint
		date        *time.Time
		severity    *entity.HeadacheSeverity
		types       *entity.HeadacheTypes
		positions   *entity.HeadachePositions
		symptoms    *entity.HeadacheSymptoms
		description *string
		isError     bool
	}{
		"update date only": {
			haID: 1,
			date: &date,
		},
		"update severity only": {
			haID:     1,
			severity: &severity,
		},
		"update everything": {
			haID:        1,
			date:        &date,
			severity:    &severity,
			types:       &types,
			positions:   &positions,
			symptoms:    &symptoms,
			description: &description,
		},
		"not found": {
			haID:    100,
			isError: true,
			date:    &date,
		},
		"invalid entry": {
			haID:      1,
			positions: &entity.HeadachePositions{"Somewhere"},
			isError:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, ctx, origHeadache := initTest(t, true)
			err := repo.PatchHeadache(ctx, test.haID, test.date, test.severity, test.types, test.positions, test.symptoms, test.description)
			if test.isError {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}
			var result entity.Headache
			repo.db.Find(&result, 1)
			expectedDate := origHeadache.Date
			if test.date != nil {
				expectedDate = date
			}
			assert.Equal(t, expectedDate, result.Date)

			expectedSeverity := origHeadache.Severity
			if test.severity != nil {
				expectedSeverity = severity
			}
			assert.Equal(t, expectedSeverity, result.Severity)

			expectedTypes := origHeadache.Types
			if test.types != nil {
				expectedTypes = types
			}
			assert.Equal(t, expectedTypes, result.Types)

			expectedPositions := origHeadache.Positions
			if test.positions != nil {
				expectedPositions = positions
			}
			assert.Equal(t, expectedPositions, result.Positions)

			expectedSymptoms := origHeadache.Symptoms
			if test.symptoms != nil {
				expectedSymptoms = symptoms
			}
			assert.Equal(t, expectedSymptoms, result.Symptoms)

			expectedDescription := origHeadache.Description
			if test.description != nil {
				expectedDescription = description
			}
			assert.Equal(t, expectedDescription, result.Description)
		})
	}
}
