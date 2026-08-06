package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	FrontendURL   string
	DbUser        string
	DbPassword    string
	DbHost        string
	DbPort        string
	DbName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	MediaDisk      string
	MediaRoot      string
	MediaURLPrefix string
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

	mediaDisk := os.Getenv("MEDIA_DISK")
	if mediaDisk == "" {
		mediaDisk = "local"
	}

	mediaRoot := os.Getenv("MEDIA_ROOT")
	if mediaRoot == "" {
		mediaRoot = "storage/media"
	}

	mediaURLPrefix := os.Getenv("MEDIA_URL_PREFIX")
	if mediaURLPrefix == "" {
		mediaURLPrefix = "/storage/media"
	}

	AppConfig = &Config{
		AppPort:     appPort,
		FrontendURL: os.Getenv("FRONTEND_URL"),
		DbUser:      os.Getenv("DB_USER"),
		DbPassword:  os.Getenv("DB_PASSWORD"),
		DbHost:      os.Getenv("DB_HOST"),
		DbPort:      os.Getenv("DB_PORT"),
		DbName:      os.Getenv("DB_NAME"),

		RedisHost:     os.Getenv("JWT_REDIS_HOST"),
		RedisPort:     os.Getenv("JWT_REDIS_PORT"),
		RedisPassword: os.Getenv("JWT_REDIS_PASSWORD"),
		RedisDB:       os.Getenv("JWT_REDIS_DB"),

		MediaDisk:      mediaDisk,
		MediaRoot:      mediaRoot,
		MediaURLPrefix: mediaURLPrefix,
	}
}
