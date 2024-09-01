package internal

import "encore.app/entity"

type SymptomsUseCase struct {
	sym ISymtpomsRepo
	cat ISymptomCategoriesRepo
}

func NewSymptomsUseCase(sym ISymtpomsRepo, cat ISymptomCategoriesRepo) SymptomsUseCase {
	return SymptomsUseCase{sym: sym, cat: cat}
}

func (uc SymptomsUseCase) List() entity.SymptomCategories {
	return uc.cat.ListCategories()
}

func (uc SymptomsUseCase) CreateCategory(name string) (uint, error) {
	return uc.cat.CreateCategory(name)
}

func (uc SymptomsUseCase) PutSymptom(symptomName string, symtpomCategoryId uint) (uint, error) {
	return uc.sym.CreateOrReplace(symptomName, symtpomCategoryId)
}
