package config

import (
	"database/sql"
	"fmt"
	"github.com/spin311/library-api/internal/repository/postgres"
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

func getConfig() *Config {
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

func SetDbs(database *sql.DB) {
	postgres.SetUserDB(database)
	postgres.SetBookDB(database)
}

func InitDatabase() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg := getConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	return db

}
