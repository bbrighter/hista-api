package notes

import (
	"time"

	"encore.dev/types/uuid"
)

type Note struct {
	ID   uint      `gorm:"primaryKey"`
	PIID uuid.UUID `gorm:"type:uuid;index"`
	Date time.Time
	Text string
}

func (n *Note) SetPiid(id uuid.UUID) {
	n.PIID = id
}
