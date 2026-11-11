package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"flag"
	"log"
	"net/http"
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
	"github.com/joho/godotenv"

	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/mail"
	"sketchdb.cozycole.net/internal/models"
)

var StaticAssets = map[string]string{
	"css": "styles.css",
	"js":  "main.js",
}

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP network address")
	debug := flag.Bool("debug", false, "debug mode")
	dev := flag.Bool("dev", false, "use dev config and set debug true")
	serveStatic := flag.Bool("serve-static", false, "serve css, js and images")

	flag.Parse()

	err := godotenv.Load()
	if err != nil && *dev {
		log.Fatal("Error loading .env file")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var dbUrl, imgBaseUrl, origin string
	var fileStorage, ArchiveStorage fileStore.FileStorageInterface
	if *dev {
		*debug = true
		*serveStatic = true

		infoLog.Println("Dev env selected, debug mode set")

		dbUrl = os.Getenv("DEV_DB_URL")
		imgBaseUrl = os.Getenv("DEV_IMG_URL")
		origin = os.Getenv("DEV_ORIGIN")

		// set paths for serving js and css
		StaticAssets["css"] = "dist/styles.css"
		StaticAssets["js"] = "dist/main.js"

		client := S3Client(
			os.Getenv("DEV_S3_ENDPOINT"),
			os.Getenv("DEV_S3_KEY"),
			os.Getenv("DEV_S3_SECRET"),
		)
		fileStorage = &fileStore.S3Storage{
			Client:     client,
			BucketName: os.Getenv("DEV_S3_BUCKET"),
		}
		archiveClient := S3Client(
			os.Getenv("DEV_S3_ARCHIVE_ENDPOINT"),
			os.Getenv("DEV_S3_ARCHIVE_KEY"),
			os.Getenv("DEV_S3_ARCHIVE_SECRET"),
		)

		ArchiveStorage = &fileStore.S3Storage{
			Client:     archiveClient,
			BucketName: os.Getenv("DEV_S3_ARCHIVE_BUCKET"),
		}

	} else {
		infoLog.Println("Production env selected")
		dbUrl = os.Getenv("DB_URL")
		imgBaseUrl = os.Getenv("IMG_URL")
		origin = os.Getenv("ORIGIN")

		err = loadAssets()
		if err != nil {
			log.Fatal("Error loading manifest found in production build")
		}
		client := S3Client(
			os.Getenv("S3_ENDPOINT"),
			os.Getenv("S3_KEY"),
			os.Getenv("S3_SECRET"),
		)
		fileStorage = &fileStore.S3Storage{
			Client:     client,
			BucketName: os.Getenv("S3_BUCKET"),
		}

		archiveClient := S3Client(
			os.Getenv("S3_ARCHIVE_ENDPOINT"),
			os.Getenv("S3_ARCHIVE_KEY"),
			os.Getenv("S3_ARCHIVE_SECRET"),
		)

		ArchiveStorage = &fileStore.S3Storage{
			Client:     archiveClient,
			BucketName: os.Getenv("S3_ARCHIVE_BUCKET"),
		}
	}

	if dbUrl == "" {
		errorLog.Fatal("Database URL not defined")
	}

	mailApiKey := os.Getenv("MAIL_API_KEY")
	mailDomain := os.Getenv("MAIL_DOMAIN")
	mailSender := os.Getenv("MAIL_SENDER")

	if mailApiKey == "" || mailDomain == "" || mailSender == "" {
		errorLog.Fatal("MAIL_API_KEY, MAIL_DOMAIN, and MAIL_SENDER must be set")
	}
	mailer := mail.New(mailDomain, mailApiKey, mailSender)

	dbpool, err := openDB(dbUrl)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbpool.Close()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = 90 * 24 * time.Hour
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	appConfig := AppConfig{
		Origin:     origin,
		DevEnv:     *dev,
		DebugMode:  *debug,
		BaseImgUrl: imgBaseUrl,

		PasswordTTL:      30 * time.Minute,
		ActivationTTL:    24 * time.Hour,
		PageSize:         24,
		MaxSearchResults: 12,
	}
	configureOAuth(origin)
	formDecoder := form.NewDecoder()
	app := &Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		FileStorage:    fileStorage,
		ArchiveStorage: ArchiveStorage,
		Categories:     &models.CategoryModel{DB: dbpool},
		Characters:     &models.CharacterModel{DB: dbpool},
		Creators:       &models.CreatorModel{DB: dbpool},
		People:         &models.PersonModel{DB: dbpool},
		Profile:        &models.ProfileModel{DB: dbpool},
		Quotes:         &models.QuoteModel{DB: dbpool},
		Recurring:      &models.RecurringModel{DB: dbpool},
		Series:         &models.SeriesModel{DB: dbpool},
		Shows:          &models.ShowModel{DB: dbpool},
		Sketches:       &models.SketchModel{DB: dbpool},
		Tags:           &models.TagModel{DB: dbpool},
		Tokens:         &models.TokenModel{DB: dbpool},
		Users:          &models.UserModel{DB: dbpool},
		Services:       NewServices(newRepositories(dbpool), fileStorage, ArchiveStorage, mailer, appConfig),
		SessionManager: sessionManager,
		Mail:           mailer,
		Assets:         StaticAssets,
		Config:         appConfig,
	}

	app.InfoLog.Println("ORIGIN: ", origin)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes("./ui/static/", *serveStatic),
	}

	infoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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

func loadAssets() error {
	f, err := os.Open("./dist/manifest.json")
	if err != nil {
		log.Println("No asset manifest found — using default asset names")
		return err
	}
	defer f.Close()

	manifest := map[string]string{}
	if err := json.NewDecoder(f).Decode(&manifest); err != nil {
		log.Printf("Failed to parse manifest.json: %v — using default asset names", err)
		return err
	}

	// Merge manifest into StaticAssets (overrides defaults if present)
	for k, v := range manifest {
		StaticAssets[k] = v
	}
	return nil
}

func init() {
	// Register the struct so encoding/gob knows how to handle it
	gob.Register(flashMessage{})
	gob.Register(PendingRegistration{})
}
