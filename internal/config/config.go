package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
	APIKey      string
	Server      string
	TAG         string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://vless_agent:supersecret@localhost:5432/vless_agent?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		APIKey:      getEnv("API_KEY", "supersecret"),
		Server:      getEnv("SERVER", "127.0.0.1:10085"),
		TAG:         getEnv("TAG", "test"),
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
