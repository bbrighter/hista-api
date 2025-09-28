package entity

import "encore.dev/types/uuid"

type ProductInstance struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name      string
	ProductId string
	Product   Product `gorm:"-"`
}

type ProductInstances []ProductInstance

type Product struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
	Apps []App  `yaml:"apps" json:"apps"`
}

type App struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

type ProductInstanceResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Product Product   `json:"product"`
}

type ProductInstancesResponse struct {
	Instances []ProductInstanceResponse `json:"instances"`
}

func (p ProductInstance) ToResponse() ProductInstanceResponse {
	return ProductInstanceResponse{
		ID:      p.ID,
		Name:    p.Name,
		Product: p.Product,
	}
}

func (ps ProductInstances) ToResponse() ProductInstancesResponse {
	var instances = []ProductInstanceResponse{}
	for _, p := range ps {
		instances = append(instances, p.ToResponse())
	}
	return ProductInstancesResponse{Instances: instances}
}
