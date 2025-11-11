package entity

import "encore.dev/types/uuid"

type Product struct {
	ID   uint      `gorm:"primaryKey;autoIncrement"`
	PIID uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex:idx_name_piid"`
	Name string    `gorm:"uniqueIndex:idx_name_piid"`
}

type Products []*Product

func (p *Product) SetPiid(piid uuid.UUID) {
	p.PIID = piid
}

type ProductResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (prod Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:   prod.ID,
		Name: prod.Name,
	}
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
}

func (prods Products) ToResponse() ProductListResponse {
	var resps = []ProductResponse{}
	for _, prod := range prods {
		resps = append(resps, prod.ToResponse())
	}
	return ProductListResponse{Products: resps}
}
