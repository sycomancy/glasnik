package infra

import (
	"os"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
	Username string
	Password string
}

type config struct {
	DB_URL       string
	TOR_URL      string
	ProxyAuth    AuthConfig
	RegistryAuth AuthConfig
	RegistryURL  string
	ProxyHost    string
	ProxyPort    string
	RegistryPort string
}

var Config config

func LoadConfig() error {
	_ = godotenv.Load() // Ignore error as .env file is optional

	Config = config{
		DB_URL:  getEnv("DB_URL", ""),
		TOR_URL: getEnv("TOR_URL", ""),
		ProxyAuth: AuthConfig{
			Username: getEnv("PROXY_USERNAME", ""),
			Password: getEnv("PROXY_PASSWORD", ""),
		},
		RegistryAuth: AuthConfig{
			Username: getEnv("REGISTRY_USERNAME", ""),
			Password: getEnv("REGISTRY_PASSWORD", ""),
		},
		RegistryURL:  getEnv("REGISTRY_URL", "http://localhost:8082"),
		ProxyHost:    getEnv("PROXY_HOST", "localhost"),
		ProxyPort:    getEnv("PROXY_PORT", "8083"),
		RegistryPort: getEnv("REGISTRY_PORT", "8082"),
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
