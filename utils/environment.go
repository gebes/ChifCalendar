package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitEnvironment(){
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetDatabaseUrl() string {
	return os.Getenv("DatabaseUrl")
}