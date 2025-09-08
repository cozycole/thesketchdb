package models

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatal(err)
	}

	teardownScript, err := os.ReadFile("../../sql/teardown.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(context.Background(), string(teardownScript))
	if err != nil {
		t.Fatal(err)
	}

	schemaScript := "../../sql/schema.sql"
	triggersScript := "../../sql/triggers.sql"
	for _, path := range []string{schemaScript, triggersScript} {
		script, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		_, err = db.Exec(context.Background(), string(teardownScript))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}

func restoreDbScript(path string) error {
	return exec.Command("psql", os.Getenv("TEST_DB_URL"), "-f", path).Run()
}

func ptr[T any](v T) *T {
	return &v
}
