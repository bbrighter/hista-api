package internal

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
)

type (
	ISymptomsRepo interface {
		CreateOrReplace(ctx context.Context, symptomName string, symptomCategoryId uint) (uint, error)
		ChangeCategory(ctx context.Context, symptomId, newCategoryId uint) error
		RenameSymptom(ctx context.Context, symptomId uint, newName string) error
	}

	ISymptomCategoriesRepo interface {
		ListCategories(ctx context.Context) ([]*entity.SymptomCategory, error)
		CreateCategory(ctx context.Context, catName string) (uint, error)
		RenameCategory(ctx context.Context, cat *entity.SymptomCategory, newName string) error
		DeleteCategory(ctx context.Context, catId uint) error
	}

	ISymptomsUseCase interface {
		PutSymptom(ctx context.Context, symptomName string, symptomCategoryId uint) (uint, error)
		List(ctx context.Context) (entity.SymptomCategories, error)
		CreateCategory(ctx context.Context, name string) (uint, error)
		ChangeCategory(ctx context.Context, symptomId, newCategoryId uint) error
		RenameSymptom(ctx context.Context, symptomId uint, newName string) error
		RenameCategory(ctx context.Context, catId uint, newName string) error
		DeleteCategory(ctx context.Context, catId uint) error
	}
)

type SymptomsUseCase struct {
	sym ISymptomsRepo
	cat ISymptomCategoriesRepo
}

func NewSymptomsUseCase(sym ISymptomsRepo, cat ISymptomCategoriesRepo) SymptomsUseCase {
	return SymptomsUseCase{sym: sym, cat: cat}
}

func (uc SymptomsUseCase) List(ctx context.Context) (entity.SymptomCategories, error) {
	return uc.cat.ListCategories(ctx)
}

func (uc SymptomsUseCase) CreateCategory(ctx context.Context, name string) (uint, error) {
	id, err := uc.cat.CreateCategory(ctx, name)
	return id, errors.MapError(err)
}

func (uc SymptomsUseCase) PutSymptom(ctx context.Context, symptomName string, symtpomCategoryId uint) (uint, error) {
	id, err := uc.sym.CreateOrReplace(ctx, symptomName, symtpomCategoryId)
	return id, errors.MapError(err)
}

func (uc SymptomsUseCase) ChangeCategory(ctx context.Context, symptomId, newCategoryId uint) error {
	return errors.MapError(uc.sym.ChangeCategory(ctx, symptomId, newCategoryId))
}
func (uc SymptomsUseCase) RenameSymptom(ctx context.Context, symptomId uint, newName string) error {
	return errors.MapError(uc.sym.RenameSymptom(ctx, symptomId, newName))
}
func (uc SymptomsUseCase) RenameCategory(ctx context.Context, catId uint, newName string) error {
	var cat = &entity.SymptomCategory{ID: catId}
	return errors.MapError(uc.cat.RenameCategory(ctx, cat, newName))
}
func (uc SymptomsUseCase) DeleteCategory(ctx context.Context, catId uint) error {
	return errors.MapError(uc.cat.DeleteCategory(ctx, catId))
}
