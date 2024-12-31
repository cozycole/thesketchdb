package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"sketchdb.cozycole.net/internal/img"
	"sketchdb.cozycole.net/internal/models"

	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	fileStorage   img.FileStorageInterface
	baseImgUrl    string
	videos        models.VideoModelInterface
	creators      models.CreatorModelInterface
	people        models.PersonModelInterface
	characters    models.CharacterModelInterface
	search        models.SearchModelInterface
	debugMode     bool
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "debug mode")
	testing := flag.Bool("testing", false, "use testing database and img storage")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var dbUrl, imgStoragePath, imgBaseUrl string
	if *testing {
		*debug = true
		infoLog.Println("Testing env selected, debug mode set")
		dbUrl = os.Getenv("TEST_DB_URL")
		imgStoragePath = os.Getenv("TEST_IMG_DISK_STORAGE")
		imgBaseUrl = os.Getenv("TEST_IMG_URL")
	} else {
		dbUrl = os.Getenv("DB_URL")
		imgStoragePath = os.Getenv("IMG_DISK_STORAGE")
		imgBaseUrl = os.Getenv("BASE_IMG_URL")
	}

	if dbUrl == "" {
		errorLog.Fatal("Database URL not defined")
	}

	if imgStoragePath == "" {
		errorLog.Fatal("Storage path not defined")
	}
	infoLog.Println(imgStoragePath)

	dbpool, err := openDB(dbUrl)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbpool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	fileStorage := img.FileStorage{Path: imgStoragePath}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		formDecoder:   formDecoder,
		fileStorage:   &fileStorage,
		videos:        &models.VideoModel{DB: dbpool},
		creators:      &models.CreatorModel{DB: dbpool},
		people:        &models.PersonModel{DB: dbpool},
		characters:    &models.CharacterModel{DB: dbpool},
		search:        &models.SearchModel{DB: dbpool},
		debugMode:     *debug,
		baseImgUrl:    imgBaseUrl,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes("./ui/static/", imgStoragePath, imgBaseUrl),
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
