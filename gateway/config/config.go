package config

import (
	"os"
)

type Config struct {
	Port                string
	ProductServiceURL   string
	OrderServiceURL     string
	NotificationService string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:              "8084",
		ProductServiceURL: "localhost:9090", // Укажи здесь порты из docker-compose
		OrderServiceURL:   "localhost:9091",
	}
	if port, exists := os.LookupEnv("Port"); exists {
		cfg.Port = port
	}
	if productServiceURL, exists := os.LookupEnv("ProductServiceURL"); exists {
		cfg.ProductServiceURL = productServiceURL
	}
	if orderServiceURL, exists := os.LookupEnv("OrderServiceURL"); exists {
		cfg.OrderServiceURL = orderServiceURL
	}
	return cfg, nil
}
