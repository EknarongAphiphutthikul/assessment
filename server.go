package main

import (
	"database/sql"
	"os"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
	"github.com/EknarongAphiphutthikul/assessment/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	log := initialLog()

	config := config.NewConfig()
	log.Info("Load Config success.")

	db := initialPostgres(config, log)
	defer db.Close()
	log.Info("Database store initial success.")
}

func initialLog() *logrus.Logger {
	return &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}

func initialPostgres(config config.Config, log *logrus.Logger) *sql.DB {
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
