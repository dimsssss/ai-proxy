package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// .env 파일 로드
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}
