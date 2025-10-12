package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPollenIntensity(t *testing.T) {
	var intensity DWDPollenIntensity

	intensity.Today = "0"
	assert.Equal(t, intensity.PollenIntensity(), NoPollen)
	intensity.Today = "0-1"
	assert.Equal(t, intensity.PollenIntensity(), NoToSmallPollen)
	intensity.Today = "1"
	assert.Equal(t, intensity.PollenIntensity(), SmallPollen)
	intensity.Today = "1-2"
	assert.Equal(t, intensity.PollenIntensity(), SmallToMediumPollen)
	intensity.Today = "2"
	assert.Equal(t, intensity.PollenIntensity(), MediumPollen)
	intensity.Today = "2-3"
	assert.Equal(t, intensity.PollenIntensity(), MediumToHighPollen)
	intensity.Today = "3"
	assert.Equal(t, intensity.PollenIntensity(), HighPollen)
}

func TestToPollen(t *testing.T) {
	var dwd = DWDPollen{
		Beifuss:  DWDPollenIntensity{Today: "0"},
		Roggen:   DWDPollenIntensity{Today: "0-1"},
		Graeser:  DWDPollenIntensity{Today: "0-1"},
		Hasel:    DWDPollenIntensity{Today: "0"},
		Esche:    DWDPollenIntensity{Today: "0"},
		Erle:     DWDPollenIntensity{Today: "0"},
		Birke:    DWDPollenIntensity{Today: "0"},
		Ambrosia: DWDPollenIntensity{Today: "2"},
	}

	pollens := dwd.ToPollen()

	assertPollenIntensity := func(t *testing.T, pollens Pollens, ty PollenType, exptected PollenIntensity) {
		var index int = -1
		for i, pol := range pollens {
			if pol.Type == ty {
				println(i, string(pol.Type), string(ty))
				index = i
				break
			}
		}
		assert.Greater(t, index, -1, string(ty))
		assert.Equal(t, exptected, pollens[index].Intensity)
	}

	assertPollenIntensity(t, pollens, Ambrosia, MediumPollen)
	assertPollenIntensity(t, pollens, Beifuss, NoPollen)
	assertPollenIntensity(t, pollens, Roggen, NoToSmallPollen)
	assertPollenIntensity(t, pollens, Esche, NoPollen)
	assertPollenIntensity(t, pollens, Erle, NoPollen)
	assertPollenIntensity(t, pollens, Birke, NoPollen)
	assertPollenIntensity(t, pollens, Graeser, NoToSmallPollen)
}
