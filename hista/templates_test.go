package hista

import (
	"encore.app/hista/entity"
	"encore.dev/beta/errs"
)

func (s *ApiTestSuite) TestListTemplates() {
	s.createTestTemplate()

	resp, err := s.service.ListTemplates(s.ctx, s.piid)
	s.NoError(err)
	s.Len(resp.Templates, 1)
	s.Len(resp.Templates[0].Items, 1)
}

func (s *ApiTestSuite) TestPostTemplate() {
	tests := map[string]struct {
		ingMissing bool
		expErrCode errs.ErrCode
	}{
		"ok":                 {},
		"ingredient missing": {ingMissing: true, expErrCode: errs.InvalidArgument},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var ingredientId uint = 1000
			if !test.ingMissing {
				_, ingredientId = s.createTestFood()
			}

			template, err := s.service.PostTemplate(
				s.ctx,
				s.piid,
				TemplateParams{
					Name: "Template",
					Items: []TemplateItemParams{
						{Condition: entity.Cooked, IngredientId: ingredientId},
					}})
			if s.assertErrCode(err, test.expErrCode) {
				return
			}
			s.NotEqualValues(0, template.ID)
		})
	}
}

func (s *ApiTestSuite) TestPutTemplate() {
	tests := map[string]struct {
		ingMissing      bool
		templateMissing bool
		expErrCode      errs.ErrCode
	}{
		"ok":                 {},
		"template missing":   {templateMissing: true, expErrCode: errs.NotFound},
		"ingredient missing": {ingMissing: true, expErrCode: errs.InvalidArgument},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var ingredientId uint = 1000
			if !test.ingMissing {
				_, ingredientId = s.createTestFood()
			}
			var templateId uint = 1000
			if !test.templateMissing {
				templateId = s.createTestTemplate()
			}

			err := s.service.PutTemplate(
				s.ctx,
				s.piid,
				templateId,
				TemplateParams{
					Name: "New name",
					Items: []TemplateItemParams{
						{Condition: entity.Raw, IngredientId: ingredientId},
					}})
			if s.assertErrCode(err, test.expErrCode) {
				return
			}

			templates, err := s.service.ListTemplates(s.ctx, s.piid)
			s.NoError(err)
			s.Len(templates.Templates, 1)
			s.Equal("New name", templates.Templates[0].Name)
			s.Len(templates.Templates[0].Items, 1)
		})
	}
}

func (s *ApiTestSuite) TestDeleteTemplate() {
	tests := map[string]struct {
		templateMissing bool
		expErrCode      errs.ErrCode
	}{
		"ok":               {},
		"template missing": {templateMissing: true, expErrCode: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var templateId uint = 1000
			if !test.templateMissing {
				templateId = s.createTestTemplate()
			}

			err := s.service.DeleteTemplate(s.ctx, s.piid, templateId)
			if s.assertErrCode(err, test.expErrCode) {
				return
			}

			resp, err := s.service.ListTemplates(s.ctx, s.piid)
			s.NoError(err)
			s.Len(resp.Templates, 0)
		})
	}
}
