package app

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	Host     string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() *Config {
	config := &Config{
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "0.0.0.0"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "avito"),
			Password: getEnv("DB_PASSWORD", "avito"),
			DBName:   getEnv("DB_NAME", "avito"),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
