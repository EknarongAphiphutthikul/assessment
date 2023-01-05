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

	err = prepareDB(config.SqlInitialize, db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func prepareDB(sql string, db DB) error {
	if sql != "" {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}
