package product_instance

import (
	"encore.app/product_mgmt/product"
	"encore.dev/types/uuid"
)

type ProductInstance struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name      string
	ProductId string
	Product   product.Product `gorm:"-"`
}
