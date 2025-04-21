package internal

import "encore.app/entity"

type SymptomsUseCase struct {
	sym ISymptomsRepo
	cat ISymptomCategoriesRepo
}

func NewSymptomsUseCase(sym ISymptomsRepo, cat ISymptomCategoriesRepo) SymptomsUseCase {
	return SymptomsUseCase{sym: sym, cat: cat}
}

func (uc SymptomsUseCase) List() entity.SymptomCategories {
	return uc.cat.ListCategories()
}

func (uc SymptomsUseCase) CreateCategory(name string) (uint, error) {
	var cat = &entity.SymptomCategory{Name: name}
	err := uc.cat.CreateCategory(cat)
	return cat.ID, err
}

func (uc SymptomsUseCase) PutSymptom(symptomName string, symtpomCategoryId uint) (uint, error) {
	return uc.sym.CreateOrReplace(symptomName, symtpomCategoryId)
}

func (uc SymptomsUseCase) ChangeCategory(symptomId, newCategoryId uint) error {
	return uc.sym.ChangeCategory(symptomId, newCategoryId)
}
func (uc SymptomsUseCase) RenameSymptom(symptomId uint, newName string) error {
	return uc.sym.RenameSymptom(symptomId, newName)
}
func (uc SymptomsUseCase) RenameCategory(catId uint, newName string) error {
	var cat = &entity.SymptomCategory{ID: catId}
	return uc.cat.RenameCategory(cat, newName)
}
func (uc SymptomsUseCase) DeleteCategory(catId uint) error {
	return uc.cat.DeleteCategory(catId)
}
