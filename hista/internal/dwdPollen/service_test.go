package dwdPollen

import (
	"context"
	"errors"
	"testing"
	"time"

	"encore.app/hista/internal/dwd"
	"encore.app/hista/internal/pollen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type mockKarlsruheDataGetter struct {
	mock.Mock
}

func (m *mockKarlsruheDataGetter) GetKarlsruheData() (dwd.DWDPollen, time.Time, error) {
	args := m.Called()
	return args.Get(0).(dwd.DWDPollen), args.Get(1).(time.Time), args.Error(2)
}

type mockPollens struct {
	mock.Mock
}

func (m *mockPollens) ListPollen(ctx context.Context) ([]pollen.PollenEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]pollen.PollenEvent), args.Error(1)
}

func (m *mockPollens) DoesExistAfter(ctx context.Context, t time.Time) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *mockPollens) CreatePollen(ctx context.Context, e *pollen.PollenEvent) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

// Tests
func TestListPollen(t *testing.T) {
	dwd := new(mockKarlsruheDataGetter)
	p := new(mockPollens)

	service := &DWDPollenService{dwd: dwd, p: p}

	p.On("ListPollen", mock.Anything).Return([]pollen.PollenEvent{pollen.PollenEvent{}}, nil)
	result, err := service.ListPollen(context.Background())

	assert.Len(t, result, 1)
	assert.NoError(t, err)
}

func TestCreatePollen_Success(t *testing.T) {
	d := new(mockKarlsruheDataGetter)
	p := new(mockPollens)

	service := &DWDPollenService{dwd: d, p: p}

	updatedAt := time.Now()
	dwdData := dwd.DWDPollen{
		Ambrosia: dwd.DWDPollenIntensity{Today: "0"},
		Roggen:   dwd.DWDPollenIntensity{Today: "1"},
		Erle:     dwd.DWDPollenIntensity{Today: "2"},
		Beifuss:  dwd.DWDPollenIntensity{Today: "0-1"},
		Birke:    dwd.DWDPollenIntensity{Today: "1-2"},
		Graeser:  dwd.DWDPollenIntensity{Today: "2-3"},
		Hasel:    dwd.DWDPollenIntensity{Today: "3"},
		Esche:    dwd.DWDPollenIntensity{Today: "0"},
	}

	tests := map[string]struct {
		isError               bool
		getKarlsruheDataError error
		doesExistAfterError   error
		createPollenError     error
	}{
		"ok":                    {isError: false},
		"error in getKarlsruhe": {isError: true, getKarlsruheDataError: errors.New("err")},
		"error in doesExists":   {isError: true, doesExistAfterError: errors.New("err")},
		"error in createPollen": {isError: true, createPollenError: errors.New("err")},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			m1 := d.On("GetKarlsruheData").Return(dwdData, updatedAt, test.getKarlsruheDataError)
			m2 := p.On("DoesExistAfter", ctx, mock.Anything).Return(test.doesExistAfterError)
			m3 := p.On("CreatePollen", ctx, mock.Anything).Return(test.createPollenError)

			err := service.CreatePollen(ctx)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			m1.Unset()
			m2.Unset()
			m3.Unset()
		})
	}
}
