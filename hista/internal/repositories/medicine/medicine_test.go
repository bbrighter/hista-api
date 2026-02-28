package medicine

import (
	"slices"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
)

func (s *MedicineRepoTestSuite) TestReorder() {
	medicines := []entity.Medicine{
		{ID: 101, PIID: GUID, Name: "Name 1", SortOrder: 100},
		{ID: 102, PIID: GUID, Name: "Name 2", SortOrder: 150},
		{ID: 103, PIID: GUID, Name: "Name 3", SortOrder: 200},
		{ID: 104, PIID: GUID, Name: "Name 4", SortOrder: 300},
		{ID: 105, PIID: GUID, Name: "Name 5", SortOrder: 250},
		{ID: 106, PIID: GUID, Name: "Name 6", SortOrder: 230},
	}
	err := s.med.db.CreateInBatches(medicines, 100).Error
	s.Require().NoError(err)

	err = s.med.Reorder(s.ctx)
	s.NoError(err)

	results, err := generic_queries.List[*entity.Medicine](s.ctx, s.med.db)
	s.Require().NoError(err)
	slices.SortFunc(results, func(a, b *entity.Medicine) int {
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
