package internalAuth

import (
	"time"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

type Token struct {
	ID      uint
	Bearer  uuid.UUID
	Expires time.Time
	UserID  uint
}

func (token Token) isValid() (bool, error) {
	if token.UserID == 0 {
		return false, errors.ErrorAttributeMustBeSet("token.UserID")
	}
	var relevantUser *User
	for _, user := range memorizedUsers {
		if user.ID == token.UserID {
			relevantUser = user
		}
	}
	for _, tok := range relevantUser.Tokens {
		if tok.Bearer == token.Bearer {
			return true, nil
		}
	}
	return false, nil
}

func cleanupTokens(service *Service) error {
	var err error
	err = service.db.Where("expires < ?", time.Now()).Delete(&Token{}).Error
	if err != nil {
		return err
	}
	updateMemorizedUsers(service)
	return nil
}

func (user User) firstOrCreateValidToken(service *Service) (Token, error) {
	var guid uuid.UUID
	var err error
	guid, err = uuid.NewV4()
	if err != nil {
		return Token{}, err
	}
	if user.ID == 0 {
		return Token{}, errors.ErrorIDMissing
	}
	var newToken = Token{
		Bearer:  guid,
		Expires: time.Now().Add(time.Hour * 24 * 7),
		UserID:  user.ID,
	}
	tx := service.db.
		Where("expires > ?", time.Now()).
		Where(Token{UserID: user.ID}).
		FirstOrCreate(&newToken)
	if tx.Error != nil {
		return Token{}, tx.Error
	}
	if tx.RowsAffected > 0 {
		if err := cleanupTokens(service); err != nil {
			return Token{}, err
		}
	}
	return newToken, nil
}
