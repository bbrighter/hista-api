package entity

import (
	"time"

	"encore.dev/types/uuid"
)

type Moment struct {
	PIID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	UpdatedAt time.Time
}

type MomentsResponse struct {
	ETag   int64          `header:"ETag"`
	ListId uint           `json:"listId"`
	Items  []ItemResponse `json:"items"`
	Status int            `encore:"httpstatus"`
	ProductListResponse
}

func (m Moment) ETag() int64 {
	return m.UpdatedAt.Unix()
}

func (m Moment) To304Response() MomentsResponse {
	return MomentsResponse{Status: 304}
}

func (m Moment) ToResponse(list *List, products Products) MomentsResponse {
	return MomentsResponse{
		ETag:                m.ETag(),
		ListId:              list.ID,
		Items:               list.ToResponse().Items,
		ProductListResponse: products.ToResponse(),
		Status:              200,
	}
}
