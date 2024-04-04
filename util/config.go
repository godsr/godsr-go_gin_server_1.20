package util

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Load() {

	if err := godotenv.Overload(".env"); err != nil {
		panic("Error loading .env file")
	}
}

func TestLoad() {

	if err := godotenv.Overload("../.env"); err != nil {
		panic("Error loading .env file")
	}
}

// Config func to get env value
func Conf(key string) string {
	var value string = os.Getenv(key)
	if strings.HasPrefix(value, "ENC+>") {
		value = value[5:]
		value = decrypt(value)
	}
	return value
}
