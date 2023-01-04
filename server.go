package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
	"github.com/EknarongAphiphutthikul/assessment/pkg/config"
	"github.com/EknarongAphiphutthikul/assessment/pkg/expenses"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log := initialLog()

	config := config.NewConfig()
	log.Info("Load Config success.")

	db := initialPostgres(config, log)
	defer db.Close()
	log.Info("Database store initial success.")

	gracefulShutdown(startServer(config, log, db), log)
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

func startServer(config config.Config, logger *logrus.Logger, db *sql.DB) *http.Server {
	handler := echoHandler(logger, db)

	srv := &http.Server{
		Addr:    ":" + config.Port(),
		Handler: handler,
	}

	logger.Infof("App started. PORT=%v", config.Port())
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("shutting down the server %v", err)
		}
	}()

	return srv
}

func echoHandler(logger *logrus.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())

	// set Expenses Handler
	expensesHandler(e, db)

	return e
}

func expensesHandler(e *echo.Echo, db *sql.DB) {
	expenDb := expenses.New(db)
	expenHandler := expenses.NewHandler(expenDb)
	g := e.Group("/expenses")
	g.POST("", expenHandler.AddExpenses)
}

func gracefulShutdown(srv *http.Server, log *logrus.Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	log.Info("App is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down: %v", err)
	}
	log.Info("Bye")
}
