package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/internal/medicines"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

type MedicineResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"isArchived"`
	SortOrder  int    `json:"sortOrder"`
}
type MedicineListResponse struct {
	Medicines []MedicineResponse `json:"medicines"`
}

func toMedicineResponse(m *medicines.Medicine) MedicineResponse {
	return MedicineResponse{ID: m.ID, Name: m.Name, IsArchived: m.IsArchived, SortOrder: m.SortOrder}
}

func toMedicineListResponse(ms medicines.Medicines) MedicineListResponse {
	medicines := []MedicineResponse{}
	for _, m := range ms {
		medicines = append(medicines, toMedicineResponse(m))
	}
	return MedicineListResponse{Medicines: medicines}
}

// encore:api auth method=GET path=/piid/:piid/medicines
func (s *Service) ListMedicines(ctx context.Context, piid uuid.UUID) (MedicineListResponse, error) {
	list, err := s.meds.ListMedicines(ctx)
	if err != nil {
		return MedicineListResponse{}, errors.MapError(err)
	}
	return toMedicineListResponse(list), nil
}

type PostMedicineParams struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/piid/:piid/medicines
func (s *Service) CreateMedicine(ctx context.Context, piid uuid.UUID, params PostMedicineParams) (IDResponse, error) {
	id, err := s.meds.CreateMedicine(ctx, params.Name)
	return IDResponse{ID: id}, errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/medicines/:medicineId
func (s *Service) DeleteMedicine(ctx context.Context, piid uuid.UUID, medicineId uint) error {
	err := s.meds.DeleteMedicine(ctx, medicineId)
	return errors.MapError(err)
}

type PatchMedicineParams struct {
	Name    *string `json:"name" encore:"optional"`
	Archive *bool   `json:"archive" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/medicines/:medicineId
func (s *Service) PatchMedicine(ctx context.Context, piid uuid.UUID, medicineId uint, params PatchMedicineParams) error {
	var err error
	if params.Name != nil {
		err = s.meds.RenameMedicine(ctx, medicineId, *params.Name)
	} else if params.Archive != nil {
		err = s.meds.ArchiveMedicine(ctx, medicineId, *params.Archive)
	} else {
		err = errors.ErrBadRequest
	}
	return errors.MapError(err)
}

type MoveMedicineParams struct {
	PreviousId *uint `json:"previousId" encore:"optional"`
	NextId     *uint `json:"nextId" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/medicines/:medicineId/reorder
func (s *Service) ReorderMedicine(ctx context.Context, piid uuid.UUID, medicineId uint, params MoveMedicineParams) error {
	if params.NextId == nil && params.PreviousId == nil {
		return errors.NewEncoreError("nextId or previousId must be set", errs.InvalidArgument)
	}
	return errors.MapError(s.meds.ReorderMedicine(ctx, medicineId, params.PreviousId, params.NextId))
}
