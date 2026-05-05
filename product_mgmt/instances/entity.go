package instances

import "encore.dev/types/uuid"

type Instance struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name      string
	ProductId string
}

func (Instance) TableName() string {
	return "product_instances"
}
