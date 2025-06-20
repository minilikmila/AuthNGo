package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}

func ComparePassword(hashPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))

	return err == nil
}
