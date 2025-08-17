package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	GoogleClientID string

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string

	TwilioSID   string
	TwilioToken string
	TwilioFrom  string

	OTPExpiryMinutes int
}

func LoadConfig() *Config {
	c := &Config{
		Port:           getEnv("PORT", "8081"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/usersvc?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "devsecret"),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),

		SMTPHost: getEnv("SMTP_HOST", ""),
		SMTPPort: getEnv("SMTP_PORT", "587"),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASS", ""),

		TwilioSID:   getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioToken: getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFrom:  getEnv("TWILIO_FROM", ""),
	}

	if v := getEnv("OTP_EXPIRY_MINUTES", "10"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.OTPExpiryMinutes = n
		} else {
			c.OTPExpiryMinutes = 10
		}
	} else {
		c.OTPExpiryMinutes = 10
	}

	return c
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
