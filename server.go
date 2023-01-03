package main

import (
	"database/sql"
	"log"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
	"github.com/EknarongAphiphutthikul/assessment/pkg/config"
)

func main() {
	config := config.NewConfig()
	db := initialPostgres(config)
	defer db.Close()
	log.Print("Database store initial success.")
}

func initialPostgres(config config.Config) *sql.DB {
	createTable := `
		CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY, title TEXT,	amount FLOAT,	note TEXT,	tags TEXT[]	);
	`
	db, err := common.NewDb(common.DbConfig{
		DriverName:    "postgres",
		Url:           config.DbUrl(),
		SqlInitialize: createTable,
	})
	if err != nil {
		log.Fatalf("Database store initial fail : %v", err)
		panic(err)
	}

	return db
}
