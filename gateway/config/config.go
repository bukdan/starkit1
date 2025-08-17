package main

import "os"

type Config struct {
	Port       string
	UserSvcURL string
}

func LoadConfig() *Config {
	c := &Config{
		Port:       getEnv("PORT_GATEWAY", "8080"),
		UserSvcURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
	}
	return c
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
