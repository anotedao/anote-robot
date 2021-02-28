package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct holds all our configuration
type Config struct {
	Dev             bool   `json:"dev"`
	Debug           bool   `json:"debug"`
	TelegramAPIKey  string `json:"telegram_api_key"`
	PostgreSQL      string `json:"postgre_sql"`
	WavesNodeAPIKey string `json:"waves_node_api_key"`
	WavesNodeHost   string `json:"waves_node_host"`
}

// Load method loads configuration file to Config struct
func (sc *Config) Load(configFile string) {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.Load] Got error while opening config file: %v", err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[Config.Load] Error while decoding JSON: %v", err)
	}
}

func initConfig() *Config {
	c := &Config{}
	c.Load("config.json")
	return c
}
