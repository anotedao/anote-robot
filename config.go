package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct holds all our configuration
type Config struct {
	TelegramAPIKey  string `json:"telegram_api_key"`
	Password        []byte `json:"password"`
	PublicKey       string `json:"public_key"`
	PrivateKey      string `json:"private_key"`
	PublicKeyToday  string `json:"public_key_today"`
	PrivateKeyToday string `json:"private_key_today"`
	Dev             bool   `json:"dev"`
	Beneficiary     string `json:"beneficiary"`
}

// Load method loads configuration file to Config struct
func (c *Config) load(configFile string) {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.load] Got error while opening config file: %v", err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)

	if err != nil {
		log.Printf("[Config.load] Error while decoding JSON: %v", err)
	}
}

func initConfig() *Config {
	c := &Config{}
	c.load("config.json")
	return c
}
