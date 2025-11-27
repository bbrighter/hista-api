package authentication

import (
	"context"

	"encore.app/shared/contextKeys"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
)

type ApiTestSuite struct {
	suite.Suite
	service *Service
	ctx     context.Context
	piid    uuid.UUID
}

func (s *ApiTestSuite) SetupSuite() {
	service, err := initService()
	s.Require().NoError(err)
	s.service = service
	s.piid = uuid.FromStringOrNil("0c5e945e-ef6c-4934-91ff-702d94e2e7a8")
	s.ctx = context.WithValue(s.T().Context(), contextKeys.Piid, s.piid)
}
