package config

import "os"

type Config struct {
	ServerAddress string
	RconAddress   string
	RconPassword  string
}

func New() *Config {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":4000"),
		RconAddress:   getEnv("RCON_ADDRESS", "127.0.0.1:25575"),
		RconPassword:  getEnv("RCON_PASSWORD", "password"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
