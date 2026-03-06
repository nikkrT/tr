package config

import "os"

type Config struct {
	Server   Server
	GRPC     GRPC
	Postgres Postgres
}

type Server struct {
	Address string
	Port    string
}

type GRPC struct {
	Port               string
	AddressGrpcProduct string
}

type Postgres struct {
	Address string
}

func load_grpc_config() GRPC {
	cfg := GRPC{
		Port:               "9091",
		AddressGrpcProduct: "localhost:9090",
	}
	if port, exists := os.LookupEnv("GRPC_PORT"); exists {
		cfg.Port = port
	}
	if address, exists := os.LookupEnv("GRPC_PRODUCT_ADDRESS"); exists {
		cfg.AddressGrpcProduct = address
	}
	return cfg
}

func loadServerConfig() Server {
	cfg := Server{
		Address: "0.0.0.0",
		Port:    "8083",
	}
	if port, exists := os.LookupEnv("ORDER_PORT"); exists {
		cfg.Port = port
	}
	if addr, exists := os.LookupEnv("ORDER_ADDRESS"); exists {
		cfg.Address = addr
	}
	return cfg
}

func loadPostgresConfig() Postgres {
	cfg := Postgres{
		Address: "postgres://postgres:password@localhost:5432/postgres",
	}
	if addr, exists := os.LookupEnv("HOST"); exists {
		cfg.Address = addr
	}
	return cfg
}

func LoadConfig() *Config {
	cfg := &Config{
		Server:   loadServerConfig(),
		GRPC:     load_grpc_config(),
		Postgres: loadPostgresConfig(),
	}
	return cfg
}
