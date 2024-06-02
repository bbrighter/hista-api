package pollen

import (
	"time"

	"gorm.io/gorm"
)

type PollenLoad uint

const (
	No            PollenLoad = 1
	NoToSmall     PollenLoad = 2
	Small         PollenLoad = 3
	SmallToMedium PollenLoad = 4
	Medium        PollenLoad = 5
	MediumToHigh  PollenLoad = 6
	High          PollenLoad = 7
)

func (load PollenLoad) String() string {
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

func (intensity PollenIntensitiy) toPollenLoad() PollenLoad {
	switch intensity.Today {
	case "0":
		return No
	case "0-1":
		return NoToSmall
	case "1":
		return Small
	case "1-2":
		return SmallToMedium
	case "2":
		return Medium
	case "2-3":
		return MediumToHigh
	default:
		return High
	}
}

type Pollen struct {
	ID        uint
	CreatedAt time.Time
	Roggen    PollenLoad
	Ambrosia  PollenLoad
	Erle      PollenLoad
	Beifuss   PollenLoad
	Birke     PollenLoad
	Graeser   PollenLoad
	Hasel     PollenLoad
	Esche     PollenLoad
}

type Pollens []Pollen

func (dwd DWD) toPollen() (Pollen, error) {
	var err error
	var karlsruhePollen DWDPollen
	karlsruhePollen, err = getKarlsruheData(dwd, Oberrhein)
	if err != nil {
		return Pollen{}, err
	}
	return Pollen{
		Roggen:   karlsruhePollen.Roggen.toPollenLoad(),
		Ambrosia: karlsruhePollen.Ambrosia.toPollenLoad(),
		Erle:     karlsruhePollen.Erle.toPollenLoad(),
		Beifuss:  karlsruhePollen.Beifuss.toPollenLoad(),
		Birke:    karlsruhePollen.Birke.toPollenLoad(),
		Graeser:  karlsruhePollen.Graeser.toPollenLoad(),
		Hasel:    karlsruhePollen.Hasel.toPollenLoad(),
		Esche:    karlsruhePollen.Esche.toPollenLoad(),
	}, nil
}

func (pollen *Pollen) writeToDatabase(service *Service, dwdLastUpdated time.Time) error {
	var mustBeUpdated bool = service.db.
		Where("created_at > ?", dwdLastUpdated).
		First(&Pollen{}).
		RowsAffected == 0
	if mustBeUpdated {
		return service.db.Create(&pollen).Error
	}
	return nil
}

func FindPollenWithSeverity(db *gorm.DB) Pollens {
	var pollens Pollens
	db.Debug().Where("roggen > 1").Or("ambrosia > 1").
		Or("erle > 1").Or("beifuss > 1").
		Or("birke > 1").Or("graeser > 1").
		Or("hasel > 1").Or("esche > 1").
		Find(&pollens)
	return pollens
}
