package config

import (
	"os"
)

type Config struct {
	Port          string
	CookieName    string
	EncryptionKey string
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8081"),
		CookieName:    getEnv("COOKIE_NAME", "cookie"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", "key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
