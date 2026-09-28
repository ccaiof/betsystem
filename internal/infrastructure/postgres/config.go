package postgres

import "os"

type Config struct {
	DSN string
}

func NewConfig() *Config {
	return &Config{
		DSN: os.Getenv("DB_URL"),
	}
}
