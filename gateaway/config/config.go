package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server Server
}

type Server struct {
	Port string
	Host string
}

func LoadConfig() (*Config, error) {
	srv := Server{}
	srv, err := LoadServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load server config: %v", err)
	}
	return &Config{srv}, nil
}

func LoadServerConfig() (Server, error) {
	srv := Server{
		Host: "localhost",
		Port: "8083",
	}
	if host, ok := os.LookupEnv("HOST"); ok {
		srv.Host = host
	}
	if port, ok := os.LookupEnv("PORT"); ok {
		srv.Port = port
	}
	return srv, nil
}
