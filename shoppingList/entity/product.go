package entity

import "encore.dev/types/uuid"

type Product struct {
	ID   uint      `gorm:"primaryKey;autoIncrement"`
	PIID uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex:idx_name_piid"`
	Name string    `gorm:"uniqueIndex:idx_name_piid"`
}

func (p *Product) SetPiid(piid uuid.UUID) {
	p.PIID = piid
}
