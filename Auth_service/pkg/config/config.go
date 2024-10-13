package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"log"
	"os"
)

type Config struct {
	AUTH_PORT string
	USER_PORT string

	DB_PORT     string
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	ACCESS_TOKEN   string
	REFRESH_TOKEN  string
	ADMIN_PASSWORD string
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{}

	config.AUTH_PORT = cast.ToString(coalesce("AUTH_PORT", ":8081"))
	config.USER_PORT = cast.ToString(coalesce("USER_PORT", ":50050"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "lp_auth"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "123321"))
	config.ACCESS_TOKEN = cast.ToString(coalesce("ACCESS_TOKEN", "JONS"))
	config.REFRESH_TOKEN = cast.ToString(coalesce("REFRESH_TOKEN", "JONS"))
	config.ADMIN_PASSWORD = cast.ToString(coalesce("ADMIN_PASSWORD", "123321"))

	return config

}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
