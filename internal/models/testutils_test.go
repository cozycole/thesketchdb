package models

import (
	"context"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		t.Fatal("TEST_DB_URL not supplied in os environment")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.New(
		"file://../../migrations",
		dbURL,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		t.Fatal(sourceErr)
	}
	if dbErr != nil {
		t.Fatal(dbErr)
	}

	t.Cleanup(func() {
		m, err := migrate.New(
			"file://../../migrations",
			dbURL,
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			t.Fatal(err)
		}

		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			t.Fatal(sourceErr)
		}
		if dbErr != nil {
			t.Fatal(dbErr)
		}

		db.Close()
	})

	return db
}

func ptr[T any](v T) *T {
	return &v
}
