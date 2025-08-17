package utils

import (
	"fmt"
)

func VerifyGoogleToken(token string) (string, error) {
	// contoh dummy verifikasi token
	if token == "" {
		return "", fmt.Errorf("token kosong")
	}
	// logika verifikasi token google di sini
	return "user@example.com", nil
}
