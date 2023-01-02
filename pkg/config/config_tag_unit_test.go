//go:build unit

package config_test

import (
	"os"
	"testing"

	"github.com/EknarongAphiphutthikul/assessment/pkg/config"
)

type ConfigEnv struct {
	Port  string
	DbUrl string
}

func setup(c ConfigEnv) func() {
	os.Setenv("DATABASE_URL", c.DbUrl)
	os.Setenv("PORT", c.Port)

	return func() {
		os.Clearenv()
	}
}

func TestConfig(t *testing.T) {
	t.Run("should return value of port and dbUrl correct when set environment PORT=2565 and DATABASE_URL=postgres://localhost:5432/postgres", func(t *testing.T) {
		wantPort := "2565"
		wantDbUrl := "postgres://localhost:5432/postgres"
		teardown := setup(ConfigEnv{
			DbUrl: wantDbUrl,
			Port:  wantPort,
		})
		defer teardown()

		cf := config.NewConfig()

		if cf.Port() != wantPort {
			t.Errorf("Port=%v; want %v", cf.Port(), wantPort)
		}
		if cf.DbUrl() != wantDbUrl {
			t.Errorf("DbUrl=%v; want %v", cf.DbUrl(), wantDbUrl)
		}

	})

	t.Run("should panic missing required environment variable: PORT when not set environment PORT", func(t *testing.T) {
		teardown := setup(ConfigEnv{
			DbUrl: "postgres://localhost:5432/postgres",
		})
		defer teardown()
		want := "missing required environment variable: PORT"

		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("The code did not panic")
			}
			if r != want {
				t.Errorf("got=%v; want %v", r, want)
			}
		}()

		config.NewConfig()
	})

	t.Run("should panic missing required environment variable: DATABASE_URL when not set environment DATABASE_URL", func(t *testing.T) {
		teardown := setup(ConfigEnv{
			Port: "2565",
		})
		defer teardown()
		want := "missing required environment variable: DATABASE_URL"

		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("The code did not panic")
			}
			if r != want {
				t.Errorf("got=%v; want %v", r, want)
			}
		}()

		config.NewConfig()
	})
}
