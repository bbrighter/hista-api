package hista

import (
	"encore.app/hista/entity"
)

func (s *ApiTestSuite) TestGetSymptoms() {
	resp, _ := s.service.ListSymptoms(s.ctx, s.piid)
	s.Len(resp.Categories, 0)

	s.createTestCondition()

	resp, _ = s.service.ListSymptoms(s.ctx, s.piid)
	s.Len(resp.Categories, 1)
	s.Equal("cat", resp.Categories[0].Name)
	s.Len(resp.Categories[0].Symptoms, 1)
	s.Equal("name", resp.Categories[0].Symptoms[0].Name)
}

func (s *ApiTestSuite) TestPostSymptomCategory() {
	resp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat"})
	s.NoError(err)
	s.Greater(resp.ID, uint(0))

}

func (s *ApiTestSuite) TestPatchSymptomName() {
	_, symptomId, _ := s.createTestCondition()

	err := s.service.PatchSymptomName(s.ctx, s.piid, symptomId, PatchSymptomNameParams{Name: "new name"})
	s.NoError(err)

	var symptom entity.Symptom
	s.service.DB.Take(&symptom, symptomId)
	s.Equal("new name", symptom.Name)
}

func (s *ApiTestSuite) TestPatchSymptomCategory() {
	_, symptomId, _ := s.createTestCondition()

	resp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "new cat"})
	s.NoError(err)

	err = s.service.PatchSymptomCategory(s.ctx, s.piid, symptomId, PatchSymptomCategoryParams{ToCategoryID: resp.ID})
	s.NoError(err)

	var symptom entity.Symptom
	s.service.DB.Take(&symptom, symptomId)
	s.Equal(resp.ID, symptom.SymptomCategoryID)
}

func (s *ApiTestSuite) TestPatchCategoryName() {
	_, _, catId := s.createTestCondition()

	err := s.service.PatchCategoryName(s.ctx, s.piid, catId, PatchCategoryNameParams{Name: "new cat name"})
	s.NoError(err)

	var cat entity.SymptomCategory
	s.service.DB.Take(&cat, catId)
	s.Equal("new cat name", cat.Name)
}

func (s *ApiTestSuite) TestDeleteSymptomCategory() {
	_, _, catId := s.createTestCondition()

	err := s.service.DeleteSymptomCategory(s.ctx, s.piid, catId)
	s.Error(err)
	rows := s.service.DB.Take(&entity.SymptomCategories{}, catId).RowsAffected
	s.EqualValues(1, rows)

	resp, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "new cat"})
	s.NoError(err)
	err = s.service.DeleteSymptomCategory(s.ctx, s.piid, resp.ID)
	s.NoError(err)
	rows = s.service.DB.Take(&entity.SymptomCategories{}, resp.ID).RowsAffected
	s.EqualValues(0, rows)
}
