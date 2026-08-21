package status

import (
	"context"
	"reflect"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"
)

type StatusService struct {
	s *statusRepo
}

func NewStatusService(db *gorm.DB) *StatusService {
	s := newStatusRepo(db)
	return &StatusService{s: s}
}

func (s *StatusService) ListStatuses(ctx context.Context) ([]*Status, error) {
	return s.s.ListStatus(ctx)
}

func (s *StatusService) CreateStatus(ctx context.Context, date time.Time) (*Status, error) {
	status := &Status{Date: date}
	err := s.s.CreateStatus(ctx, status)
	return status, err
}

func (s *StatusService) DeleteStatus(ctx context.Context, id uint) error {
	return s.s.DeleteStatus(ctx, id)
}

type UpdateStatusParams struct {
	Date                  *time.Time
	MorningFitness        *int
	MorningSleep          *int
	DayFitness            *int
	EveningFitness        *int
	Depressive            *int
	Tense                 *int
	MoodSwings            *int
	Irritable             *int
	LossOfInterest        *int
	ConcentrationProblems *int
	LackOfDrive           *int
	AppetiteChanges       *int
	SleepProblems         *int
	Overwhelmed           *int
	Crash                 *bool
}

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func (p UpdateStatusParams) ToMap() map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(p)
	typ := reflect.TypeOf(p)
	for i := range val.NumField() {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldVal.IsNil() {
			continue
		}
		dbName := toSnakeCase(field.Name)
		result[dbName] = fieldVal.Elem().Interface()
	}
	return result
}

func (s *StatusService) UpdateStatus(
	ctx context.Context,
	id uint,
	params UpdateStatusParams,
) error {
	values := params.ToMap()
	return s.s.UpdateStatus(ctx, id, values)
}
