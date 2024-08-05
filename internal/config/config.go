package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Address string
	Storage
}

type Storage struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	return &Config{
		os.Getenv("ADDRESS"),
		Storage{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			DB:       os.Getenv("POSTGRES_DB"),
		},
	}
}
