package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func GetString(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}

	return fallback
}
