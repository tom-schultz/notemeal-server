package internal

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
)

func CreateSecureString() string {
	return rand.Text()
}

func CompareHashAndString(hash string, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}

func HashString(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	return string(hash), err
}
