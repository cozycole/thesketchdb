package main

import (
	"context"
	"flag"
	"html/template"
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
	"sketchdb.cozycole.net/internal/img"
	"sketchdb.cozycole.net/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	templateCache  map[string]*template.Template
	fileStorage    img.FileStorageInterface
	baseImgUrl     string
	cast           models.CastModelInterface
	categories     models.CategoryInterface
	characters     models.CharacterModelInterface
	creators       models.CreatorModelInterface
	people         models.PersonModelInterface
	profile        models.ProfileModelInterface
	shows          models.ShowModelInterface
	tags           models.TagModelInterface
	users          models.UserModelInterface
	sketches       models.SketchModelInterface
	sessionManager *scs.SessionManager
	debugMode      bool
	formDecoder    *form.Decoder
	settings       settings
}

type settings struct {
	pageSize          int
	maxSearchResults  int
	localImageServer  bool
	localImageStorage bool
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "debug mode")
	testing := flag.Bool("testing", false, "use testing config and set debug true")
	localImgServer := flag.Bool("localimg", false, "serve images from local directory")
	localImgStorage := flag.Bool("localstorage", false, "store/delete images in local directory")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var dbUrl, imgStoragePath, imgBaseUrl string
	var fileStorage img.FileStorageInterface
	if *testing {
		*debug = true
		infoLog.Println("Testing env selected, debug mode set")
		dbUrl = os.Getenv("TEST_DB_URL")
		imgBaseUrl = os.Getenv("TEST_IMG_URL")
		if *localImgServer {
			imgStoragePath = os.Getenv("TEST_IMG_DISK_STORAGE")
			fileStorage = &img.FileStorage{RootPath: imgStoragePath}
		} else {
			client := S3Client(
				os.Getenv("TEST_S3_ENDPOINT"),
				os.Getenv("TEST_S3_KEY"),
				os.Getenv("TEST_S3_SECRET"),
			)
			fileStorage = &img.S3Storage{
				Client:     client,
				BucketName: os.Getenv("TEST_S3_BUCKET"),
			}
		}
	} else {
		infoLog.Println("Production env selected")
		dbUrl = os.Getenv("DB_URL")
		imgBaseUrl = os.Getenv("IMG_URL")
		client := S3Client(
			os.Getenv("S3_ENDPOINT"),
			os.Getenv("S3_KEY"),
			os.Getenv("S3_SECRET"),
		)
		fileStorage = &img.S3Storage{
			Client:     client,
			BucketName: os.Getenv("S3_BUCKET"),
		}
	}

	if dbUrl == "" {
		errorLog.Fatal("Database URL not defined")
	}

	if imgStoragePath == "" && *localImgStorage {
		errorLog.Fatal("Storage path not defined")
	}

	dbpool, err := openDB(dbUrl)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbpool.Close()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		fileStorage:    fileStorage,
		cast:           &models.CastModel{DB: dbpool},
		categories:     &models.CategoryModel{DB: dbpool},
		characters:     &models.CharacterModel{DB: dbpool},
		creators:       &models.CreatorModel{DB: dbpool},
		people:         &models.PersonModel{DB: dbpool},
		profile:        &models.ProfileModel{DB: dbpool},
		shows:          &models.ShowModel{DB: dbpool},
		tags:           &models.TagModel{DB: dbpool},
		users:          &models.UserModel{DB: dbpool},
		sketches:       &models.SketchModel{DB: dbpool},
		sessionManager: sessionManager,
		debugMode:      *debug,
		baseImgUrl:     imgBaseUrl,
		settings: settings{
			pageSize:          24,
			maxSearchResults:  12,
			localImageServer:  *localImgServer,
			localImageStorage: *localImgStorage,
		},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes("./ui/static/", imgStoragePath),
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
