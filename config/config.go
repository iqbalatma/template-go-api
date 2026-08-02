package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	FrontendURL string
	DbUser      string
	DbPassword  string
	DbHost      string
	DbPort      string
	DbName      string
}

var AppConfig *Config

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8000"
	}
	AppConfig = &Config{
		AppPort:     appPort,
		FrontendURL: os.Getenv("FRONTEND_URL"),
		DbUser:      os.Getenv("DB_USER"),
		DbPassword:  os.Getenv("DB_PASSWORD"),
		DbHost:      os.Getenv("DB_HOST"),
		DbPort:      os.Getenv("DB_PORT"),
		DbName:      os.Getenv("DB_NAME"),
	}
}
