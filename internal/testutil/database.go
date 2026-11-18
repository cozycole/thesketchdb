package testutil

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTestDb(
	t *testing.T,
	adminDSN string,
	migrationsDir string,
) *pgxpool.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to the administrative database.
	admin, err := pgxpool.New(ctx, adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()

	if err := admin.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	// Generate a unique database name.
	dbName := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	// Database identifiers cannot be SQL parameters.
	// This is safe because dbName is generated internally.
	_, err = admin.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err != nil {
		t.Fatal(err)
	}

	// Register cleanup immediately after creating the database.
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()

		admin, err := pgxpool.New(cleanupCtx, adminDSN)
		if err != nil {
			t.Errorf("connect for cleanup: %v", err)
			return
		}
		defer admin.Close()

		// Close any remaining connections to the test database.
		_, err = admin.Exec(
			cleanupCtx,
			fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, dbName),
		)
		if err != nil {
			t.Errorf("drop test database %s: %v", dbName, err)
		}
	})

	// Build a connection string targeting the new database.
	u, err := url.Parse(adminDSN)
	if err != nil {
		t.Fatal(err)
	}

	u.Path = "/" + dbName
	testDSN := u.String()

	// Apply migrations.
	if err := ApplyMigrations(testDSN, "file://"+migrationsDir); err != nil {
		t.Fatal(err)
	}

	// Open application connection pool.
	db, err := pgxpool.New(ctx, testDSN)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		t.Fatal(err)
	}

	// Registered after DROP cleanup, so it executes first.
	t.Cleanup(db.Close)

	return db
}
