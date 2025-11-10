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
	Quantity *uint8

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (i *Item) SetPiid(piid uuid.UUID) {
	i.PIID = piid
	i.ProductPiid = piid
	i.ListPiid = piid
}

type ItemResponse struct {
	ID        uint      `json:"id"`
	ProductId uint      `json:"productId"`
	ListId    uint      `json:"listId"`
	Checked   bool      `json:"checked"`
	Quantity  *uint8    `json:"quantity,omitempty" encore:"optional"`
	CreatedAt time.Time `json:"-"`
}

func (item Item) ToResponse() ItemResponse {
	return ItemResponse{
		ID:        item.ID,
		ListId:    item.ListId,
		ProductId: item.ProductId,
		Checked:   item.Checked,
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
	}
}
