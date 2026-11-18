package integration_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/mail"
	"github.com/cozycole/thesketchdb/internal/testutil"
	"github.com/cozycole/thesketchdb/internal/web"
)

// TestApp contains the resources needed by integration tests.
type TestApp struct {
	App    *web.Application
	Server *httptest.Server
	DB     *pgxpool.Pool
	Client *http.Client
}

// TestOptions allows individual tests to supply dependencies.
type TestOptions struct {
	Mailer         mail.MailerInterface
	FileStorage    fileStore.FileStorageInterface
	ArchiveStorage fileStore.FileStorageInterface
}

// newTestDB establishes a connection to the integration test database.
// newTestSessionManager configures sessions for HTTP integration tests.
func newTestSessionManager(db *pgxpool.Pool) *scs.SessionManager {
	sessionManager := scs.New()

	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 90 * 24 * time.Hour

	// httptest.Server uses HTTP rather than HTTPS.
	sessionManager.Cookie.Secure = false

	return sessionManager
}

// newTestConfig provides application configuration for tests.
func newTestConfig() web.AppConfig {
	return web.AppConfig{
		Origin:           "http://localhost",
		DevEnv:           true,
		DebugMode:        true,
		PasswordTTL:      30 * time.Minute,
		ActivationTTL:    24 * time.Hour,
		PageSize:         24,
		MaxSearchResults: 12,
	}
}

// newTestDependencies constructs the application's dependencies.
func newTestDependencies(
	t *testing.T,
	db *pgxpool.Pool,
	opts TestOptions,
) web.AppDependencies {
	t.Helper()
	templateCache, err := web.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	return web.AppDependencies{
		DB:             db,
		SessionManager: newTestSessionManager(db),
		FormDecoder:    form.NewDecoder(),

		Mailer:         opts.Mailer,
		FileStorage:    opts.FileStorage,
		ArchiveStorage: opts.ArchiveStorage,

		InfoLog:  log.New(io.Discard, "", 0),
		ErrorLog: log.New(io.Discard, "", 0),

		Assets: map[string]string{
			"css": "styles.css",
			"js":  "main.js",
		},

		TemplateCache: templateCache,
	}
}

// newTestApp initializes the application and starts an HTTP server.
func newTestApp(
	t *testing.T,
	opts TestOptions,
) *TestApp {
	t.Helper()

	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		t.Fatal("TEST_DB_URL not defined")
	}
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		t.Fatal("MIGRATIONS_PATH not defined")
	}

	db := testutil.NewTestDb(t, dbURL, migrationsPath)
	config := newTestConfig()
	deps := newTestDependencies(t, db, opts)

	app := web.NewApplication(config, deps)

	server := httptest.NewServer(app.Routes("", false))
	t.Cleanup(server.Close)

	return &TestApp{
		App:    app,
		Server: server,
		DB:     db,
		Client: server.Client(),
	}
}
