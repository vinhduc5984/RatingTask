package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
		log.Println("Warning: No .env file found. Continuing...")
	}

	return os.Getenv("MONGOURI")
}
