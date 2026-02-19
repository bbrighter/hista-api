package hista

func (s *ApiTestSuite) TestMedicine() {
	listResp, err := s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 0)

	idResp, err := s.service.CreateMedicine(s.ctx, s.piid, PostMedicineParams{Name: "Name"})
	s.NoError(err)
	medicineId := idResp.ID

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 1)
	s.Equal(listResp.Medicines[0].Name, "Name")
	s.Equal(listResp.Medicines[0].IsArchived, false)

	var newName string = "new name"
	err = s.service.PatchMedicine(s.ctx, s.piid, medicineId, PatchMedicineParams{Name: &newName})
	s.NoError(err)

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 1)
	s.Equal(listResp.Medicines[0].Name, "new name")

	var archive bool = true
	err = s.service.PatchMedicine(s.ctx, s.piid, medicineId, PatchMedicineParams{Archive: &archive})
	s.NoError(err)

	listResp, err = s.service.ListMedicines(s.ctx, s.piid)
	s.NoError(err)
	s.Len(listResp.Medicines, 1)
	s.Equal(listResp.Medicines[0].IsArchived, true)
}
