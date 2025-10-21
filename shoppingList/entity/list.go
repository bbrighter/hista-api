package entity

import (
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type List struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	PIID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	DeletedAt gorm.DeletedAt
	Items     []Item `gorm:"foreignKey:ListId,ListPiid;references:ID,PIID"`
}

func (l *List) SetPiid(piid uuid.UUID) {
	l.PIID = piid
}

type ListResponse struct {
	ID    uint           `json:"id"`
	Items []ItemResponse `json:"items"`
}

func (list List) ToResponse() ListResponse {
	var itemsResp = []ItemResponse{}
	for _, item := range list.Items {
		itemsResp = append(itemsResp, item.ToResponse())
	}
	return ListResponse{
		ID:    list.ID,
		Items: itemsResp,
	}
}
