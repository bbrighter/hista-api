package errors

import (
	"context"

	"encore.dev/rlog"
	"encore.dev/types/option"
)

// encore:service
type Service struct{}

func initService() (*Service, error) {
	return &Service{}, nil
}

type ErrorParams struct {
	Text    string                `json:"text"`
	Status  option.Option[int]    `json:"status"`
	Details option.Option[string] `json:"details"`
	Stack   option.Option[string] `json:"stack"`
}

// encore:api auth method=POST path=/error
func (s *Service) LogError(ctx context.Context, params ErrorParams) error {
	rlog.Error("frontend error logged",
		"Text", params.Text,
		"Status", params.Status.GetOrElse(0),
		"Details", params.Details.GetOrElse(""),
		"Stack", params.Stack.GetOrElse(""),
	)
	return nil
}
