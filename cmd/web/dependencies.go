package main

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/mail"
)

type AppDependencies struct {
	DB             *pgxpool.Pool
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
	FormDecoder    *form.Decoder
	Mailer         mail.MailerInterface
	FileStorage    fileStore.FileStorageInterface
	ArchiveStorage fileStore.FileStorageInterface
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	Assets         map[string]string // depends on ./dist/manifest.json
}

func newDependencies(
	config AppConfig,
) (AppDependencies, func(), error) {
	dev := config.DevEnv

	dbpool, err := openDB(env("DB_URL", dev))
	if err != nil {
		return AppDependencies{}, nil, err
	}

	cleanup := func() {
		dbpool.Close()
	}

	success := false
	defer func() {
		if !success {
			cleanup()
		}
	}()

	client := S3Client(
		env("S3_ENDPOINT", dev),
		env("S3_KEY", dev),
		env("S3_SECRET", dev),
	)
	fileStorage := &fileStore.S3Storage{
		Client:     client,
		BucketName: env("S3_BUCKET", dev),
	}

	archiveClient := S3Client(
		env("S3_ARCHIVE_ENDPOINT", dev),
		env("S3_ARCHIVE_KEY", dev),
		env("S3_ARCHIVE_SECRET", dev),
	)
	archiveStorage := &fileStore.S3Storage{
		Client:     archiveClient,
		BucketName: env("S3_ARCHIVE_BUCKET", dev),
	}

	mailApiKey := os.Getenv("MAIL_API_KEY")
	mailDomain := os.Getenv("MAIL_DOMAIN")
	mailSender := os.Getenv("MAIL_SENDER")
	if mailApiKey == "" || mailDomain == "" || mailSender == "" {
		return AppDependencies{}, nil, errors.New("error: mailer config")
	}

	mailer := mail.New(mailDomain, mailApiKey, mailSender)

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = 90 * 24 * time.Hour
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()
	if err != nil {
		return AppDependencies{}, nil, err
	}

	formDecoder := form.NewDecoder()

	var assets = map[string]string{
		"css": "dist/styles.css",
		"js":  "dist/main.js",
	}

	if !dev {
		assets, err = loadAssets(assets)
		if err != nil {
			return AppDependencies{}, nil, err
		}
	}

	success = true
	return AppDependencies{
		DB:             dbpool,
		FileStorage:    fileStorage,
		ArchiveStorage: archiveStorage,
		Mailer:         mailer,
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
		InfoLog:        log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:       log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Assets:         assets,
	}, cleanup, nil
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = dbpool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return dbpool, nil
}

func S3Client(endpoint, key, secret string) *s3.S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession := session.Must(session.NewSession(s3Config))
	return s3.New(newSession)
}

func loadAssets(staticAssets map[string]string) (map[string]string, error) {
	f, err := os.Open("./dist/manifest.json")
	if err != nil {
		log.Println("No asset manifest found — using default asset names")
		return nil, err
	}
	defer f.Close()

	manifest := map[string]string{}
	if err := json.NewDecoder(f).Decode(&manifest); err != nil {
		log.Printf("Failed to parse manifest.json: %v — using default asset names", err)
		return nil, err
	}

	// Merge manifest into StaticAssets (overrides defaults)
	for k, v := range manifest {
		staticAssets[k] = v
	}
	return staticAssets, nil
}

func env(key string, dev bool) string {
	if dev {
		return os.Getenv("DEV_" + key)
	}
	return os.Getenv(key)
}
