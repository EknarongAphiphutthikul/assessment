package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	DriverName    string
	Url           string
	SqlInitialize string
}

type DbStorage struct {
	db *sql.DB
}

func NewDb(config DbConfig) (*DbStorage, error) {
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

	return &DbStorage{db: db}, nil
}

func (s *DbStorage) TearDown() {
	s.db.Close()
	log.Printf("DB Storage Teardown.")
}
