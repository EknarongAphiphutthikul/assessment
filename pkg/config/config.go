package config

import "os"

func getenv(name string, require bool, defaultVAlue string) string {
	v := os.Getenv(name)
	if v == "" {
		if require {
			panic("missing required environment variable: " + name)
		}
		return defaultVAlue
	}
	return v
}

type Config struct {
	port    string
	dbUrl   string
	authKey string
}

func NewConfig() Config {
	return Config{
		port:    getenv("PORT", true, ""),
		dbUrl:   getenv("DATABASE_URL", true, ""),
		authKey: getenv("AUTH_KEY", false, "November 10, 2009"),
	}
}

func (c Config) Port() string {
	return c.port
}

func (c Config) DbUrl() string {
	return c.dbUrl
}

func (c Config) AuthKey() string {
	return c.authKey
}
