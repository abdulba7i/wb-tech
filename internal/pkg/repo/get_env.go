package repo

import (
	"os"
	"strings"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func HidePassword(connStr string) string {
	const passwordKey = "password="
	start := strings.Index(connStr, passwordKey)
	if start == -1 {
		return connStr
	}
	start += len(passwordKey)
	end := strings.Index(connStr[start:], " ")
	if end == -1 {
		return connStr[:start] + "***"
	}
	return connStr[:start] + "***" + connStr[start+end:]
}
