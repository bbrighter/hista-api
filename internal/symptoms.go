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
