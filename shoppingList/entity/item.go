package entity

import (
	"time"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID   uint      `gorm:"primaryKey;autoIncrement"`
	PIID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ProductId   uint
	ProductPiid uuid.UUID `gorm:"type:uuid"`
	Product     Product   `gorm:"foreignKey:ProductId,ProductPiid;references:ID,PIID"`

	ListId   uint
	ListPiid uuid.UUID `gorm:"type:uuid"`
	List     List      `gorm:"foreignKey:ListId,ListPiid;references:ID,PIID"`

	Checked  bool
	Quantity uint8

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (i *Item) SetPiid(piid uuid.UUID) {
	i.PIID = piid
	i.ProductPiid = piid
	i.ListPiid = piid
}
