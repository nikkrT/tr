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
		serverPort: "8080",
		serverHost: "localhost",
	}
	if port, exists := os.LookupEnv("PORT"); exists {
		cfg.serverPort = port
	}
	if host, exists := os.LookupEnv("HOST"); exists {
		cfg.serverHost = host
	}
	return cfg
}

func LoadConfigRabbitMQ() RabbitMQ {
	cfg := RabbitMQ{
		addr:           "amqp://guest:guest@localhost:5672/",
		Exchange:       "direct",
		Queue:          "products",
		RoutingKey:     "products",
		ConsumerTag:    "products",
		WorkerPoolSize: 1,
	}
	if addr, exists := os.LookupEnv("HOST"); exists {
		cfg.addr = addr
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
		Addr: "postgres://postgres:postgres@localhost:5432/products",
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
