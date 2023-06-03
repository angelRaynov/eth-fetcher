package config

import (
	"log"
	"os"
)

type Application struct {
	AppMode    string
	APIPort    string
	APIKey     string
	EthNodeURL string

	DSN string
}

func New() *Application {
	var config Application

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		log.Fatal("API_PORT not set")
	}
	config.APIPort = apiPort

	ethURL := os.Getenv("ETH_NODE_URL")
	if ethURL == "" {
		log.Fatal("ETH_NODE_URL not set")
	}
	config.EthNodeURL = ethURL

	dsn := os.Getenv("DB_CONNECTION_URL")
	if dsn == "" {
		log.Fatal("DB_CONNECTION_URL not set")
	}
	config.DSN = dsn

	return &config
}
