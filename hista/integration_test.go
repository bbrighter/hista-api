package hista

import (
	"time"

	"encore.app/hista/entity"
	"encore.dev/beta/errs"
)

func (s *ApiTestSuite) TestHeadaches() {
	var err error
	var headaches entity.HeadachesResponse
	headaches, err = s.service.ListHeadaches(s.ctx, s.piid)
	s.NoError(err)
	s.Len(headaches.Headaches, 0)

	// Create one headache
	date := time.Date(2018, 1, 2, 3, 4, 5, 0, time.Local)
	var severity entity.HeadacheSeverity = 5
	idResp, err := s.service.PostHeadache(s.ctx, s.piid, PostHeadacheParams{Date: date, Severity: severity})
	id := idResp.ID
	s.NoError(err)
	s.Greater(id, uint(0))

	// Get headache
	headaches, err = s.service.ListHeadaches(s.ctx, s.piid)
	s.NoError(err)
	s.Len(headaches.Headaches, 1)

	headache, err := s.service.GetHeadache(s.ctx, s.piid, id)
	s.NoError(err)
	s.Equal(date, headache.Date)
	s.Equal(severity, headache.Severity)
	s.Nil(headache.Positions)

	// Patch and verify
	newDate := time.Date(2019, 1, 2, 3, 4, 5, 0, time.Local)
	err = s.service.PatchHeadacheDate(s.ctx, s.piid, id, PatchHeadacheDateParams{Date: newDate})
	s.NoError(err)
	var newSeverity entity.HeadacheSeverity = 1
	err = s.service.PatchHeadacheSeverity(s.ctx, s.piid, id, PatchHeadacheSeverityParams{Severity: newSeverity})
	s.NoError(err)
	newTypes := entity.HeadacheTypes{entity.Dull}
	err = s.service.PatchHeadacheTypes(s.ctx, s.piid, id, PatchHeadacheTypesParams{Types: newTypes})
	s.NoError(err)
	newPositions := entity.HeadachePositions{entity.Back, entity.Ear}
	err = s.service.PatchHeadachePositions(s.ctx, s.piid, id, PatchHeadachePositionsParams{Positions: newPositions})
	s.NoError(err)
	newSymptoms := entity.HeadacheSymptoms{entity.Dizziness}
	err = s.service.PatchHeadacheSymptoms(s.ctx, s.piid, id, PatchHeadacheSymptomsParams{Symptoms: newSymptoms})
	s.NoError(err)

	headache, err = s.service.GetHeadache(s.ctx, s.piid, id)
	s.NoError(err)
	s.Equal(newDate, headache.Date)
	s.Equal(newPositions, headache.Positions)
	s.Equal(newSeverity, headache.Severity)
	s.Equal(newSymptoms, headache.Symptoms)
	s.Equal(newTypes, headache.Types)

	// Delete and verify
	err = s.service.DeleteHeadache(s.ctx, s.piid, idResp.ID)
	s.NoError(err)
	headache, err = s.service.GetHeadache(s.ctx, s.piid, idResp.ID)
	s.Error(err)
	headaches, err = s.service.ListHeadaches(s.ctx, s.piid)
	s.NoError(err)
	s.Len(headaches.Headaches, 0)
}

func (s *ApiTestSuite) TestSymptoms() {
	symptoms, err := s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 0)

	// Create a category
	catId, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat"})
	s.NoError(err)

	// Create a symptom
	event, err := s.service.CreateConditionEvent(s.ctx, s.piid)
	s.NoError(err)
	symptomName := "symptom"
	_, err = s.service.PostCondition(s.ctx, s.piid, event.ID, ConditionRequestParams{SymptomName: &symptomName, CategoryID: &catId.ID})
	s.NoError(err)

	symptoms, err = s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 1)
	cat := symptoms.Categories[0]
	s.Equal("cat", cat.Name)
	s.Len(cat.Symptoms, 1)
	sym := cat.Symptoms[0]
	s.Equal("symptom", sym.Name)

	// Patch category and symptom names
	err = s.service.PatchCategoryName(s.ctx, s.piid, cat.ID, PatchCategoryNameParams{Name: "new cat"})
	s.NoError(err)
	err = s.service.PatchSymptomName(s.ctx, s.piid, sym.ID, PatchSymptomNameParams{Name: "new symptom"})
	s.NoError(err)
	symptoms, err = s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 1)
	cat = symptoms.Categories[0]
	s.Equal("new cat", cat.Name)
	s.Len(cat.Symptoms, 1)
	sym = cat.Symptoms[0]
	s.Equal("new symptom", sym.Name)

	// Change category
	cat2Id, err := s.service.PostSymptomCategory(s.ctx, s.piid, PostSymptomCategoryRequest{Name: "cat2"})
	s.NoError(err)
	err = s.service.PatchSymptomCategory(s.ctx, s.piid, sym.ID, PatchSymptomCategoryParams{ToCategoryID: cat2Id.ID})
	s.NoError(err)
	symptoms, err = s.service.ListSymptoms(s.ctx, s.piid)
	s.NoError(err)
	s.Len(symptoms.Categories, 2)
	for _, c := range symptoms.Categories {
		if c.ID == cat2Id.ID {
			s.Equal("new symptom", c.Symptoms[0].Name)
		}
	}

	// Delete
	err = s.service.DeleteSymptomCategory(s.ctx, s.piid, cat.ID)
	s.NoError(err)
	err = s.service.DeleteSymptomCategory(s.ctx, s.piid, cat2Id.ID)
	s.assertErrCode(err, errs.InvalidArgument)
}
