package internal

import (
	"context"
	"testing"

	"encore.app/users/entity"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
)

type testUserRepo struct{}

const GUID_STR = "dd8a9474-6443-45aa-b81a-c561c887e0a4"

func (r testUserRepo) Create(ctx context.Context, name string, password string, instanceId uuid.UUID) (entity.User, error) {
	return entity.User{}, nil
}
func (r testUserRepo) Find(ctx context.Context, id uuid.UUID) (entity.User, error) {
	guid, _ := uuid.FromString(GUID_STR)
	return entity.User{ID: guid, Name: "Name", Password: "$2y$12$8qHSy4rs4o/dxlUOjgoiO.unuMenyjpBHtMJi55VyHipTXGsPF0uS"}, nil
}
func (r testUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (r testUserRepo) ChangePassword(ctx context.Context, id uuid.UUID, password string) error {
	return nil
}

func TestChangePassword(t *testing.T) {
	uc := NewUserUseCase(testUserRepo{})
	ctx := context.Background()

	guid, _ := uuid.FromString(GUID_STR)
	err := uc.ChangePassword(ctx, guid, "new password", "Password")
	assert.NoError(t, err)

	err = uc.ChangePassword(ctx, guid, "Very new password", "Wrong password")
	assert.Error(t, err)
}
