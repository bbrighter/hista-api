package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/beta/errs"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/medicines
func (s *Service) ListMedicines(ctx context.Context, piid uuid.UUID) (entity.MedicineListResponse, error) {
	list, err := s.medicineList.List(ctx)
	return list.ToResponse(), errors.MapError(err)
}

type PostMedicineParams struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/piid/:piid/medicines
func (s *Service) CreateMedicine(ctx context.Context, piid uuid.UUID, params PostMedicineParams) (entity.IDResponse, error) {
	id, err := s.medicine.Create(ctx, params.Name)
	return entity.IDResponse{ID: id}, errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/medicines/:medicineId
func (s *Service) DeleteMedicine(ctx context.Context, piid uuid.UUID, medicineId uint) error {
	return errors.MapError(s.medicine.Delete(ctx, medicineId))
}

type PatchMedicineParams struct {
	Name    *string `json:"name" encore:"optional"`
	Archive *bool   `json:"archive" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/medicines/:medicineId
func (s *Service) PatchMedicine(ctx context.Context, piid uuid.UUID, medicineId uint, params PatchMedicineParams) error {
	var err error
	if params.Name != nil {
		err = s.medicine.Rename(ctx, medicineId, *params.Name)
	} else if params.Archive != nil {
		err = s.medicine.Archive(ctx, medicineId, *params.Archive)
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
	return errors.MapError(s.medicine.Reorder(ctx, medicineId, params.PreviousId, params.NextId))
}
