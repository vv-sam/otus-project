package services

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Users struct {
	users map[string]string
}

func NewUsers(pairs map[string]string) *Users {
	return &Users{users: pairs}
}

func (u *Users) Login(username, password string) bool {
	return u.users[username] == password
}

func (u *Users) GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	return token.SignedString([]byte("random-secret-key"))
}

func (u *Users) ValidateToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("random-secret-key"), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["username"].(string), nil
	}

	return "", errors.New("invalid token")
}
