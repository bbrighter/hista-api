package templates

import (
	"encore.app/hista/entity"
	"gorm.io/gorm"
)

func (s *TemplateTestSuite) TestList() {
	templates, err := s.s.List(s.ctx)

	s.NoError(err)
	s.Len(templates, 0)

	s.createTemplate()
	templates, err = s.s.List(s.ctx)
	s.NoError(err)
	s.Len(templates, 1)
	s.Len(templates[0].Items, 2)
}

func (s *TemplateTestSuite) TestCreateEmpty() {
	id, err := s.s.Create(s.ctx, "empty template", []entity.TemplateItem{})
	s.NoError(err)

	template := s.getTemplateById(id)
	s.Equal("empty template", template.Name)
	s.Len(template.Items, 0)
}

func (s *TemplateTestSuite) TestCreate() {
	s.createIngredients()

	id, err := s.s.Create(s.ctx, "filled template", []entity.TemplateItem{
		{Condition: entity.Cooked, IngredientID: s.ingIds[0]},
		{Condition: entity.Cooked, IngredientID: s.ingIds[1]},
	})
	s.NoError(err)

	template := s.getTemplateById(id)
	s.Equal("filled template", template.Name)
	s.Len(template.Items, 2)
}

func (s *TemplateTestSuite) TestUpdate() {
	s.createTemplate()

	ingId := s.ingIds[1]
	err := s.s.Update(s.ctx, s.templateID, "new name", []entity.TemplateItem{{Condition: entity.Cooked, IngredientID: ingId}})
	s.NoError(err)

	template := s.getTemplateById(s.templateID)
	s.Equal("new name", template.Name)
	s.Len(template.Items, 1)
	s.Equal(entity.Cooked, template.Items[0].Condition)
	s.Equal(ingId, template.Items[0].IngredientID)
}

func (s *TemplateTestSuite) TestUpdateTemplateNotFound() {
	s.createTemplate()

	ingId := s.ingIds[1]
	err := s.s.Update(s.ctx, 1000, "new name", []entity.TemplateItem{{Condition: entity.Cooked, IngredientID: ingId}})
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *TemplateTestSuite) TestUpdateIngredientNotFound() {
	s.createTemplate()

	err := s.s.Update(s.ctx, s.templateID, "new name", []entity.TemplateItem{{Condition: entity.Cooked, IngredientID: 1000}})
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *TemplateTestSuite) TestDelete() {
	s.createTemplate()

	err := s.s.Delete(s.ctx, s.templateID)
	s.NoError(err)

	noTemplates := count[*entity.Template](s)
	s.EqualValues(0, noTemplates)
	noItems := count[*entity.TemplateItem](s)
	s.EqualValues(0, noItems)
}

func (s *TemplateTestSuite) TestDeleteNotFound() {
	err := s.s.Delete(s.ctx, 1000)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *TemplateTestSuite) TestGet() {
	s.createTemplate()

	template, err := s.s.Get(s.ctx, s.templateID)
	s.NoError(err)
	s.Equal("template", template.Name)
	s.Len(template.Items, 2)
}

func (s *TemplateTestSuite) TestGetNotFound() {
	_, err := s.s.Get(s.ctx, 1000)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}
