package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/internal/symptoms"
	"encore.dev/types/uuid"
)

type SymptomResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CategoryID uint   `json:"categoryId"`
}

type SymptomCategoryResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Symptoms []SymptomResponse `json:"symptoms"`
}

type SymptomCategoryListResponse struct {
	Categories []SymptomCategoryResponse
}

func toSymptomResponse(s symptoms.Symptom) SymptomResponse {
	return SymptomResponse{
		ID:         s.ID,
		Name:       s.Name,
		CategoryID: s.SymptomCategoryID,
	}
}

func toSymptomCategoryListResponse(cats symptoms.SymptomCategories) SymptomCategoryListResponse {
	var resp = []SymptomCategoryResponse{}
	for _, cat := range cats {
		var symptomsResp = []SymptomResponse{}
		for _, sym := range cat.Symptoms {
			symptomsResp = append(symptomsResp, toSymptomResponse(sym))
		}
		var c = SymptomCategoryResponse{
			ID:       cat.ID,
			Name:     cat.Name,
			Symptoms: symptomsResp,
		}
		resp = append(resp, c)
	}
	return SymptomCategoryListResponse{Categories: resp}
}

// encore:api auth method=GET path=/piid/:piid/symptoms
func (service *Service) ListSymptoms(ctx context.Context, piid uuid.UUID) (SymptomCategoryListResponse, error) {
	categories, err := service.syms.ListSymptomCategories(ctx)
	if err != nil {
		return SymptomCategoryListResponse{}, errors.MapError(err)
	}
	return toSymptomCategoryListResponse(categories), nil
}

type PatchSymptomNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/piid/:piid/symptoms/:id/name
func (service *Service) PatchSymptomName(ctx context.Context, piid uuid.UUID, id uint, params PatchSymptomNameParams) error {
	return errors.MapError(service.syms.RenameSymptom(ctx, id, params.Name))
}

type PatchSymptomCategoryParams struct {
	ToCategoryID uint `json:"toCategoryId"`
}

// encore:api auth method=PATCH path=/piid/:piid/symptoms/:id/category
func (service *Service) PatchSymptomCategory(ctx context.Context, piid uuid.UUID, id uint, params PatchSymptomCategoryParams) error {
	return errors.MapError(service.syms.ChangeCategory(ctx, id, params.ToCategoryID))
}

type PostSymptomCategoryRequest struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/piid/:piid/symptom-categories
func (service *Service) PostSymptomCategory(ctx context.Context, piid uuid.UUID, params PostSymptomCategoryRequest) (IDResponse, error) {
	id, err := service.syms.CreateCategory(ctx, params.Name)
	return IDResponse{ID: id}, errors.MapError(err)
}

type PatchCategoryNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/piid/:piid/symptom-categories/:id
func (service *Service) PatchCategoryName(ctx context.Context, piid uuid.UUID, id uint, params PatchCategoryNameParams) error {
	return errors.MapError(service.syms.RenameCategory(ctx, id, params.Name))
}

// encore:api auth method=DELETE path=/piid/:piid/symptom-categories/:id
func (service *Service) DeleteSymptomCategory(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(service.syms.DeleteCategory(ctx, id))
}
