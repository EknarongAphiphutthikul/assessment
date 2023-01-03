package common

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	DriverName    string
	Url           string
	SqlInitialize string
}

func NewDb(config DbConfig) (*sql.DB, error) {
	db, err := sql.Open(config.DriverName, config.Url)
	if err != nil {
		return nil, err
	}

	if config.SqlInitialize != "" {
		_, err = db.Exec(config.SqlInitialize)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
