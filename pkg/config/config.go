package config

import "os"

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable: " + name)
	}
	return v
}

type Config struct {
	port  string
	dbUrl string
}

func NewConfig() Config {
	return Config{
		port:  getenv("PORT"),
		dbUrl: getenv("DATABASE_URL"),
	}
}

func (c Config) Port() string {
	return c.port
}

func (c Config) DbUrl() string {
	return c.dbUrl
}
