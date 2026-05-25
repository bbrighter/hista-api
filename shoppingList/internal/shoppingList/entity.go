package shoppinglist

import (
	"time"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID   uint      `gorm:"primaryKey;autoIncrement"`
	PIID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ProductId   uint      `gorm:"uniqueIndex:idx_product_list"`
	ProductPiid uuid.UUID `gorm:"type:uuid"`
	Product     Product   `gorm:"foreignKey:ProductId,ProductPiid;references:ID,PIID"`

	ListId   uint      `gorm:"uniqueIndex:idx_product_list"`
	ListPiid uuid.UUID `gorm:"type:uuid"`
	List     List      `gorm:"foreignKey:ListId,ListPiid;references:ID,PIID"`

	Checked  bool
	Quantity *uint8

	CreatedAt time.Time
}

type Product struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	PIID     uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex:idx_name_piid"`
	Name     string    `gorm:"uniqueIndex:idx_name_piid"`
	Archived bool      `gorm:"not null;default:false"`
}

type Products []*Product

func (p *Product) SetPiid(piid uuid.UUID) {
	p.PIID = piid
}

func (i *Item) SetPiid(piid uuid.UUID) {
	i.PIID = piid
	i.ProductPiid = piid
	i.ListPiid = piid
}

type List struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	PIID      uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex:idx_list_singleton,where:deleted_at IS NULL"`
	DeletedAt gorm.DeletedAt
	Items     []Item `gorm:"foreignKey:ListId,ListPiid;references:ID,PIID"`
	Singleton bool   `gorm:"not null;default:true;uniqueIndex:idx_list_singleton,where:deleted_at IS NULL"`
}

func (l *List) SetPiid(piid uuid.UUID) {
	l.PIID = piid
}
