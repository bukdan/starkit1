package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword wraps bcrypt.GenerateFromPassword
func HashPassword(password string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// ComparePassword checks hash vs plain password
func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
