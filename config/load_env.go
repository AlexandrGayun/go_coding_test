package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv(path ...string) {
	err := godotenv.Load(path...)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}
