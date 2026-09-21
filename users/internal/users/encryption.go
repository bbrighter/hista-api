package users

import "golang.org/x/crypto/bcrypt"

type Encryption struct {
	costs int
}

func NewEncryption(costs int) *Encryption {
	return &Encryption{costs: costs}
}

func (e *Encryption) GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), e.costs)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (e *Encryption) ValidatePassword(password1, password2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
}
