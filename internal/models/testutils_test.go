package models

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	godotenv.Load()
	db, err := pgxpool.New(context.Background(), os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatal(err)
	}

	schemaDirPath := "../../sql/schema"
	entries, err := os.ReadDir(schemaDirPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		scriptPath := path.Join(schemaDirPath, entry.Name())
		script, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("../../sql/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})
	return db
}
