package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
	Config   DBConfig
}

type DBConfig struct {
	DSN string
}

func LoadConfig() (Config, error) {
	config := Config{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.HTTPAddr = os.Getenv("HTTP_ADDR")
	if config.HTTPAddr == "" {
		log.Println("HTTP_ADDR env variable not set, setting default to http://0.0.0.0:8080")
		config.HTTPAddr = ":8080"
	}
	config.Config.DSN = os.Getenv("DATABASE_URL")
	if config.Config.DSN == "" {
		log.Println("DATABASE_URL env variable not set")
		return Config{}, errors.New("DATABASE_URL env variable not set")
	}
	return config, nil
}
