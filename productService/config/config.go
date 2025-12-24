package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQ
	Postgres PostgresConfig
}

type ServerConfig struct {
	ServerPort string
	ServerHost string
}

type RabbitMQ struct {
	Addr           string
	Exchange       string
	Queue          string
	RoutingKey     string
	ConsumerTag    string
	WorkerPoolSize int
}

type PostgresConfig struct {
	Addr string
}

func LoadConfigServer() ServerConfig {
	cfg := ServerConfig{
		ServerPort: "8080",
		ServerHost: "localhost",
	}
	if port, exists := os.LookupEnv("PORT"); exists {
		cfg.ServerPort = port
	}
	if host, exists := os.LookupEnv("HOST"); exists {
		cfg.ServerHost = host
	}
	return cfg
}

func LoadConfigRabbitMQ() RabbitMQ {
	cfg := RabbitMQ{
		Addr:           "amqp://guest:guest@localhost:5672/",
		Exchange:       "productService",
		Queue:          "products",
		RoutingKey:     "products",
		ConsumerTag:    "products",
		WorkerPoolSize: 1,
	}
	if addr, exists := os.LookupEnv("HOST"); exists {
		cfg.Addr = addr
	}
	if Exchange, exists := os.LookupEnv("EXCHANGE_NAME"); exists {
		cfg.Exchange = Exchange
	}
	if Queue, exists := os.LookupEnv("QUEUE_NAME"); exists {
		cfg.Queue = Queue
	}
	if ConsumerTag, exists := os.LookupEnv("CONSUMER_TAG_NAME"); exists {
		cfg.ConsumerTag = ConsumerTag
	}
	return cfg
}

func LoadConfigPostgres() PostgresConfig {
	cfg := PostgresConfig{
		Addr: "postgres://postgres:password@localhost:5432/postgres",
	}
	if addr, exists := os.LookupEnv("HOST"); exists {
		cfg.Addr = addr
	}
	return cfg
}

func LoadConfig() *Config {
	cfg := &Config{
		Server:   LoadConfigServer(),
		RabbitMQ: LoadConfigRabbitMQ(),
		Postgres: LoadConfigPostgres(),
	}
	return cfg
}
