package config

import "os"

type Config struct {
	Server Server
	Amqp   RabbitMQ
}

type Server struct {
	Host string
	Port string
}

type RabbitMQ struct {
	Addr           string
	Exchange       string
	QueueCreate    string
	QueueUpdate    string
	QueueDelete    string
	RoutingKey     string
	ConsumerTag    string
	WorkerPoolSize int
	Keys           []string
}

func LoadConfig() *Config {
	return &Config{
		Server: loadServer(),
		Amqp:   loadRabbitmq(),
	}
}
func loadRabbitmq() RabbitMQ {
	keys := []string{"product.created",
		"product.updated",
		"product.deleted"}
	cfg := RabbitMQ{
		Addr:           "amqp://guest:guest@localhost:5672/",
		QueueCreate:    "product.created",
		QueueUpdate:    "product.updated",
		QueueDelete:    "product.deleted",
		WorkerPoolSize: 3,
		Keys:           keys,
		Exchange:       "productService",
	}
	if addr, err := os.LookupEnv("RABBIT_MQ_ADDRESS"); err {
		cfg.Addr = addr
	}
	return cfg
}

func loadServer() Server {
	cfg := Server{
		Host: "localhost",
		Port: "8081",
	}
	if host, err := os.LookupEnv("SERVER_HOST"); err {
		cfg.Host = host
	}
	if port, err := os.LookupEnv("SERVER_PORT"); err {
		cfg.Port = port
	}
	return cfg
}
