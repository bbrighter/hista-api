package moments

import (
	"encore.dev/types/uuid"
)

type Moment struct {
	PIID            uuid.UUID
	ProductsVersion int
	ItemsVersion    int
}
