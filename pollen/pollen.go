package pollen

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PollenIntensity uint

const (
	No            PollenIntensity = 1
	NoToSmall     PollenIntensity = 2
	Small         PollenIntensity = 3
	SmallToMedium PollenIntensity = 4
	Medium        PollenIntensity = 5
	MediumToHigh  PollenIntensity = 6
	High          PollenIntensity = 7
)

type PollenType string

const (
	Roggen   PollenType = "Roggen"
	Ambrosia PollenType = "Ambrosia"
	Erle     PollenType = "Erle"
	Beifuss  PollenType = "Beifuss"
	Birke    PollenType = "Birke"
	Graeser  PollenType = "Gräser"
	Hasel    PollenType = "Hasel"
	Esche    PollenType = "Esche"
)

func (load PollenIntensity) String() string {
	var loadString string
	switch load {
	case No:
		loadString = "Keine"
	case NoToSmall:
		loadString = "Keine bis geringe"
	case Small:
		loadString = "Geringe"
	case SmallToMedium:
		loadString = "Geringe bis mittlere"
	case Medium:
		loadString = "Mittlere"
	case MediumToHigh:
		loadString = "Mittlere bis hohe"
	case High:
		loadString = "Hohe"
	}
	return loadString
}

type PollenEvent struct {
	ID        uint
	CreatedAt time.Time
	Pollens   Pollens
}

type Pollen struct {
	ID            uint
	PollenEventID uint
	Type          PollenType
	Intensity     PollenIntensity
}

type Pollens []Pollen

func (dwd DWD) Pollens() (Pollens, error) {
	var err error
	var karlsruhePollen DWDPollen
	karlsruhePollen, err = dwd.getKarlsruheData(Oberrhein)
	if err != nil {
		return Pollens{}, err
	}
	return Pollens{
		Pollen{Type: Ambrosia, Intensity: karlsruhePollen.Ambrosia.PollenIntensity()},
		Pollen{Type: Roggen, Intensity: karlsruhePollen.Roggen.PollenIntensity()},
		Pollen{Type: Erle, Intensity: karlsruhePollen.Erle.PollenIntensity()},
		Pollen{Type: Beifuss, Intensity: karlsruhePollen.Beifuss.PollenIntensity()},
		Pollen{Type: Birke, Intensity: karlsruhePollen.Birke.PollenIntensity()},
		Pollen{Type: Graeser, Intensity: karlsruhePollen.Graeser.PollenIntensity()},
		Pollen{Type: Hasel, Intensity: karlsruhePollen.Hasel.PollenIntensity()},
		Pollen{Type: Esche, Intensity: karlsruhePollen.Esche.PollenIntensity()},
	}, nil
}

func (pollenEvent *PollenEvent) create(service *Service, dwdLastUpdated time.Time) error {
	var mustBeUpdated bool = service.db.
		Where("created_at > ?", dwdLastUpdated).
		First(&PollenEvent{}).
		RowsAffected == 0
	if mustBeUpdated {
		return service.db.Create(&pollenEvent).Error
	}
	return nil
}

func FindPollenWithSeverity(db *gorm.DB) []PollenEvent {
	var pollenEvent []PollenEvent
	db.Preload(clause.Associations, "intensity > 1").Find(&pollenEvent)
	return pollenEvent
}

func (pollenEvent *PollenEvent) BeforeDelete(tx *gorm.DB) error {
	return tx.Delete(&Pollen{}, &Pollen{PollenEventID: pollenEvent.ID}).Error
}
