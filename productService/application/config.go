package application

import (
	"os"
)

type Config struct {
	serverPort string
	serverHost string
	dbAddr     string
}

func LoadConfig() Config {
	cfg := Config{
		serverPort: "8080",
		serverHost: "localhost",
		dbAddr:     "postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
	}
	if port, exists := os.LookupEnv("PORT"); exists {
		cfg.serverPort = port
	}
	if host, exists := os.LookupEnv("HOST"); exists {
		cfg.serverHost = host
	}
	if dbPort, exists := os.LookupEnv("DB_PORT"); exists {
		cfg.dbAddr = dbPort
	}
	return cfg
}
