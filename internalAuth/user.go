package internalAuth

import (
	"crypto/sha256"
	"encoding/hex"

	"encore.app/errors"
	"gorm.io/gorm/clause"
)

type User struct {
	ID              uint
	Name            string
	PasswordHashHex string
	Tokens          []Token
}

type Users []User

func (service *Service) initializeUsers() error {
	var userJulia = User{Name: "Julia", PasswordHashHex: "4ca3e38ab18a9c49bac7ebd3dd0067db640c630e1dd3848dd3e6f8ba351ada99"}
	return service.db.FirstOrCreate(&User{}, userJulia).Error
}

func updateMemorizedUsers(service *Service) {
	service.db.Preload(clause.Associations).Find(&memorizedUsers)
}

func hashing(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hexHash := hex.EncodeToString(hasher.Sum(nil))
	return hexHash
}

func (user User) isValidPassword(password string) (bool, error) {
	if user.Name == "" {
		return false, errors.ErrorAttributeMustBeSet("user.Name")
	}
	var relevantUser *User
	var err error
	relevantUser, err = getUserByName(user.Name)
	if err != nil {
		return false, err
	}
	var isValid bool = relevantUser.PasswordHashHex == hashing(password)
	return isValid, nil
}

func getUserByName(userName string) (*User, error) {
	for _, user := range memorizedUsers {
		if user.Name == userName {
			return user, nil
		}
	}
	return nil, errors.ErrorNotFound
}
