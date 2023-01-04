package config

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type Instance struct {
	DB     *sql.DB
	Log    *logrus.Logger
	Config *Config
}
