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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func getConfig() *Config {

	return &Config{
		DbHost:     os.Getenv("DBHOST"),
		DbPort:     os.Getenv("DBPORT"),
		DbUser:     os.Getenv("DBUSER"),
		DbPassword: os.Getenv("DBPASSWORD"),
		DbName:     os.Getenv("DBNAME"),
	}
}

func SetDbs(database *sql.DB) {
	postgres.SetUserDB(database)
	postgres.SetBookDB(database)
}

func InitDatabase() (*sql.DB, error) {
	cfg := getConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database!")

	return db, nil
}

func GetEnvString(key string) string {
	return os.Getenv(key)
}
