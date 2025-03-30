package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	Port       string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../config/.env")
	if err != nil {
		return nil, err
	}
	return &Config{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		Port:       os.Getenv("PORT"),
	}, nil
}
