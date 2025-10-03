package authentication

import (
	"encore.app/authentication/internal"
	"encore.app/authentication/internal/repository"
)

// encore:service
type Service struct {
	g internal.ITokenGenerator
}

func initService() (*Service, error) {
	keys := newKeys()
	r := repository.NewTokenRepo(keys.privateKey, keys.publicKey)
	g := internal.NewTokenGenerator(r)

	return &Service{g: g}, nil
}
