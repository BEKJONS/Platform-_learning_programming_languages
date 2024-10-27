package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PRODUCTS_PORT string

	DB_PORT     string
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

func Load() *Config {
	config := &Config{}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.PRODUCTS_PORT = os.Getenv("PRODUCTS_PORT")
	config.DB_PORT = os.Getenv("DB_PORT")
	config.DB_HOST = os.Getenv("DB_HOST")
	config.DB_USER = os.Getenv("DB_USER")
	config.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	config.DB_NAME = os.Getenv("DB_NAME")

	return config
}
