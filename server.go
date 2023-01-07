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

	cf := config.NewConfig()
	log.Info("Load Config success.")

	db := initialPostgres(cf, log)
	defer db.Close()
	log.Info("Database store initial success.")

	ins := &config.Instance{
		DB:     db,
		Log:    log,
		Config: &cf,
	}

	gracefulShutdown(startServer(ins), log)
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

func initialPostgres(config config.Config, log common.Log) *sql.DB {
	createTable := `
		CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY, title TEXT,	amount FLOAT,	note TEXT,	tags TEXT[]	);
	`
	db, err := common.NewDb(common.DbConfig{
		DriverName:    "postgres",
		Url:           config.DbUrl(),
		SqlInitialize: createTable,
	})
	if err != nil {
		log.Fatalf("Database store initial fail : %s", err)
		panic(err)
	}

	return db
}

func startServer(ins *config.Instance) *http.Server {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	initMiddleware(e, ins)
	initRoutes(e, ins)

	srv := &http.Server{
		Addr:    ":" + ins.Config.Port(),
		Handler: e,
	}

	ins.Log.Infof("App started. PORT=%s", ins.Config.Port())
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ins.Log.Fatalf("shutting down the server %s", err)
		}
	}()

	return srv
}

func gracefulShutdown(srv *http.Server, log common.Log) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	log.Info("App is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down: %s", err)
	}
	log.Info("Bye")
}

func initMiddleware(e *echo.Echo, ins *config.Instance) {
	/*
		e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
			Handler: func(c echo.Context, req []byte, resp []byte) {
				ins.Log.Info("XXX")
			},
		}))
	*/
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogMethod:   true,
		LogHeaders:  []string{"Content-Type", "Authorization"},
		LogLatency:  true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			if values.Error == nil {
				ins.Log.WithFields(logrus.Fields{
					"URI":     values.URI,
					"status":  values.Status,
					"method":  values.Method,
					"headers": values.Headers,
					"latency": values.Latency,
				}).Info("request")
			} else {
				ins.Log.WithFields(logrus.Fields{
					"URI":     values.URI,
					"status":  values.Status,
					"method":  values.Method,
					"headers": values.Headers,
					"latency": values.Latency,
					"error":   values.Error,
				}).Error("request error")
			}
			return nil
		},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			value := c.Request().Header.Values("Authorization")
			if value != nil && value[0] == ins.Config.AuthKey() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}
	})
	e.Use(middleware.Recover())
}

func initRoutes(echo *echo.Echo, ins *config.Instance) {
	expenses.Routes(echo, ins)
}
