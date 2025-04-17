package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func Load() (*Config, error) {
	// Debug: Print all environment variables

	port, err := strconv.Atoi(getRequiredEnv("SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(getRequiredEnv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			Host:     getRequiredEnv("DB_HOST"),
			Port:     dbPort,
			User:     getRequiredEnv("DB_USER"),
			Password: getRequiredEnv("DB_PASSWORD"),
			DBName:   getRequiredEnv("DB_NAME"),
		},
	}, nil
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set (value was empty)", key))
	}
	return value
}
