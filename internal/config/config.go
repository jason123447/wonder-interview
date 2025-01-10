package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SECRET_KEY string

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SECRET_KEY = os.Getenv("SECRET_KEY")
}
