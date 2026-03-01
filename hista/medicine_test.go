package hista

import (
	"slices"

	"encore.app/hista/entity"
)

func (s *ApiTestSuite) TestMedicine() {
	listResp, err := s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 0)

	idResp, err := s.service.CreateMedicine(s.ctx, s.piid, PostMedicineParams{Name: "Name"})
	s.NoError(err)
	medicine1Id := idResp.ID

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 1)
	s.Equal(listResp.Medicines[0].Name, "Name")
	s.Equal(listResp.Medicines[0].IsArchived, false)

	var newName string = "new name"
	err = s.service.PatchMedicine(s.ctx, s.piid, medicine1Id, PatchMedicineParams{Name: &newName})
	s.NoError(err)

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 1)
	s.Equal(listResp.Medicines[0].Name, "new name")

	var archive bool = true
	err = s.service.PatchMedicine(s.ctx, s.piid, medicine1Id, PatchMedicineParams{Archive: &archive})
	s.NoError(err)

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 1)
	s.Equal(listResp.Medicines[0].IsArchived, true)

	idResp2, err := s.service.CreateMedicine(s.ctx, s.piid, PostMedicineParams{Name: "Name 2"})
	s.NoError(err)
	medicine2Id := idResp2.ID
	// Now, order is Med2, Med1
	idResp3, err := s.service.CreateMedicine(s.ctx, s.piid, PostMedicineParams{Name: "Name 3"})
	s.NoError(err)
	medicine3Id := idResp3.ID
	// Now, order is Med3, Med2, Med1
	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 3)
	slices.SortFunc(listResp.Medicines, func(a, b entity.MedicineResponse) int {
		return a.SortOrder - b.SortOrder
	})
	s.Equal(100, listResp.Medicines[0].SortOrder)
	s.Equal(200, listResp.Medicines[1].SortOrder)
	s.Equal(300, listResp.Medicines[2].SortOrder)

	s.service.ReorderMedicine(
		s.ctx, s.piid,
		medicine2Id,
		MoveMedicineParams{NextId: &medicine3Id},
	) // Now, order is Med2, Med3, Med1

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 3)
	slices.SortFunc(listResp.Medicines, func(a, b entity.MedicineResponse) int {
		return a.SortOrder - b.SortOrder
	})
	s.Equal(medicine2Id, listResp.Medicines[0].ID)
	s.Equal(medicine3Id, listResp.Medicines[1].ID)
	s.Equal(medicine1Id, listResp.Medicines[2].ID)

	s.service.ReorderMedicine(
		s.ctx, s.piid,
		medicine1Id,
		MoveMedicineParams{NextId: &medicine3Id, PreviousId: &medicine2Id},
	) // Now, order is Med2, Med1, Med3

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 3)
	slices.SortFunc(listResp.Medicines, func(a, b entity.MedicineResponse) int {
		return a.SortOrder - b.SortOrder
	})
	s.Equal(medicine2Id, listResp.Medicines[0].ID)
	s.Equal(medicine1Id, listResp.Medicines[1].ID)
	s.Equal(medicine3Id, listResp.Medicines[2].ID)
}
