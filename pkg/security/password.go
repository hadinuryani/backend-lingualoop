package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password tidak boleh kosong")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
