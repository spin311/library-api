package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		DbHost:     os.Getenv("HOST"),
		DbPort:     os.Getenv("PORT"),
		DbUser:     os.Getenv("USER"),
		DbPassword: os.Getenv("PASSWORD"),
		DbName:     os.Getenv("DBNAME"),
	}
}
