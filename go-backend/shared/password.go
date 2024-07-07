package shared

import (
	"os"

	"github.com/joho/godotenv"
)

var Password string

func LoadPassword() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	Password = os.Getenv("PASSWORD")

	if Password == "" {
		panic("No password found in .env file")
	}
}
