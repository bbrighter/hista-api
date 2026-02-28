package entity

import (
	"time"

	"encore.dev/types/uuid"
)

type Medicine struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	PIID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"uniqueIndex"`
	SortOrder  int       `gorm:"not null"`
	IsArchived bool      `gorm:"not null"`
	Intakes    []Intake
}

func (m *Medicine) SetPiid(id uuid.UUID) {
	m.PIID = id
}

type Medicines []*Medicine

type Intake struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	PIID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date         time.Time `gorm:"not null"`
	Medicine     Medicine  `gorm:"foreignKey:MedicineID,MedicinePIID;referencec:ID,PIID"`
	MedicineID   uint
	MedicinePIID uuid.UUID
}

func (i *Intake) SetPiid(id uuid.UUID) {
	i.PIID = id
	i.MedicinePIID = id
}

type MedicineResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"isArchived"`
	SortOrder  int    `json:"sortOrder"`
}
type MedicineListResponse struct {
	Medicines []MedicineResponse `json:"medicines"`
}

func (m Medicine) ToResponse() MedicineResponse {
	return MedicineResponse{ID: m.ID, Name: m.Name, IsArchived: m.IsArchived, SortOrder: m.SortOrder}
}

func (ms Medicines) ToResponse() MedicineListResponse {
	medicines := []MedicineResponse{}
	for _, m := range ms {
		medicines = append(medicines, m.ToResponse())
	}
	return MedicineListResponse{Medicines: medicines}
}

type GroupedIntake struct {
	Date       time.Time `gorm:"column:date"`
	MedicineId uint      `gorm:"column:medicine_id"`
	Count      int64     `gorm:"column:count"`
}

type GroupedIntakeList []GroupedIntake

func (gil GroupedIntakeList) ToResponse() IntakeResponseList {
	intakes := []IntakeResponse{}
	for _, gi := range gil {
		intakes = append(intakes,
			IntakeResponse{
				Date:       gi.Date,
				MedicineId: gi.MedicineId,
				Count:      gi.Count,
			})
	}
	return IntakeResponseList{Intakes: intakes}
}

type IntakeResponse struct {
	Date       time.Time `json:"date"`
	MedicineId uint      `json:"medicineId"`
	Count      int64     `json:"count"`
}

type IntakeResponseList struct {
	Intakes []IntakeResponse `json:"intakes"`
}
