package config

import (
	"fmt"
	"os"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbSSLMode  string
	DebugLevel string
}

var config *Config

func getConfig() *Config {
	if config == nil {
		config = LoadConfig()
	}
	return config
}

func LoadConfig() *Config {
	return &Config{
		DbHost:     GetEnv("DB_HOST", "localhost"),
		DbPort:     GetEnv("DB_PORT", "5432"),
		DbUser:     GetEnv("DB_USER", "postgres"),
		DbPassword: GetEnv("DB_PASSWORD", "password"),
		DbName:     GetEnv("DB_NAME", "digital_wallet"),
		DbSSLMode:  GetEnv("DB_SSLMODE", "disable"),
		DebugLevel: GetEnv("DEBUG_LEVEL", "info"),
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (c *Config) GetDbConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.DbHost, c.DbPort, c.DbUser, c.DbName, c.DbPassword, c.DbSSLMode)
}
