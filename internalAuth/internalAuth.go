package internalAuth

import (
	"math/rand"
	"strings"
	"time"
)

type UserId string

const (
	Julia UserId = "Julia"
	Benni UserId = "Benni"
)

type AuthParams struct {
	UserId   UserId
	Password string
}

var users *[]AuthParams = new([]AuthParams)

func InitUsers() {
	var userJulia = AuthParams{
		UserId:   Julia,
		Password: "123pi",
	}
	var userBenni = AuthParams{
		UserId:   Benni,
		Password: "123pi",
	}
	users = &[]AuthParams{
		userBenni,
		userJulia,
	}
}

type Token struct {
	UserId  UserId
	Bearer  string
	Expires time.Time
}

var CurrentToken *Token

func (token *Token) InitToken() *Token {
	if token == nil {
		token = new(Token)
	}
	sb := strings.Builder{}
	sb.Grow(50)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 50; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	var tok string = sb.String()
	var userId UserId = Julia
	var expires time.Time = time.Now().Add(time.Hour)
	token.Bearer = tok
	token.UserId = userId
	token.Expires = expires
	return token
}

func (token *Token) isValidToken() bool {
	return token.UserId == Julia && token.Bearer == CurrentToken.Bearer
}

func IsValidCredential(password string, userName UserId) bool {
	return password == "123pi" && userName == Julia
}
