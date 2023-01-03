package main

import (
	"log"

	"github.com/EknarongAphiphutthikul/assessment/pkg/config"
	"github.com/EknarongAphiphutthikul/assessment/pkg/db"
)

func main() {
	config := config.NewConfig()

	createTable := `
		CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY, title TEXT,	amount FLOAT,	note TEXT,	tags TEXT[]	);
	`
	dbStore, err := db.NewDb(db.DbConfig{
		DriverName:    "postgres",
		Url:           config.DbUrl(),
		SqlInitialize: createTable,
	})
	if err != nil {
		log.Fatalf("Database store initial fail : %v", err)
		panic(err)
	}
	defer dbStore.TearDown()
	log.Print("Database store initial success.")
}
