package internalAuth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitUsers(t *testing.T) {
	InitUsers()

	assert.Len(t, *users, 2)
	var names []string
	var passwords []string
	for _, user := range *users {
		names = append(names, string(user.UserId))
		passwords = append(passwords, user.Password)
	}
	assert.Contains(t, names, "Julia")
	assert.Contains(t, names, "Benni")
	assert.Contains(t, passwords, "123pi")
}
