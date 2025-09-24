package main

import (
	"context"
	"encoding/json"
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
	moments        models.MomentModelInterface
	people         models.PersonModelInterface
	profile        models.ProfileModelInterface
	recurring      models.RecurringModelInterface
	shows          models.ShowModelInterface
	series         models.SeriesModelInterface
	tags           models.TagModelInterface
	users          models.UserModelInterface
	sketches       models.SketchModelInterface
	sessionManager *scs.SessionManager
	debugMode      bool
	formDecoder    *form.Decoder
	assets         map[string]string
	settings       settings
}

type settings struct {
	pageSize          int
	maxSearchResults  int
	localImageServer  bool
	localImageStorage bool
	origin            string
}

var StaticAssets = map[string]string{
	"css": "styles.css",
	"js":  "main.js",
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "HTTP network address")
	debug := flag.Bool("debug", false, "debug mode")
	dev := flag.Bool("dev", false, "use dev config and set debug true")
	localImgServer := flag.Bool("localimg", false, "serve images from local directory")
	localImgStorage := flag.Bool("localstorage", false, "store/delete images in local directory")
	serveStatic := flag.Bool("serve-static", false, "serve css, js and images")

	flag.Parse()

	err := godotenv.Load()
	if err != nil && *dev {
		log.Fatal("Error loading .env file")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var dbUrl, imgStoragePath, imgBaseUrl, origin string
	var fileStorage img.FileStorageInterface
	if *dev {
		*debug = true
		*serveStatic = true

		infoLog.Println("Testing env selected, debug mode set")

		dbUrl = os.Getenv("DEV_DB_URL")
		imgBaseUrl = os.Getenv("DEV_IMG_URL")
		origin = os.Getenv("DEV_ORIGIN")

		// set paths for serving js and css
		StaticAssets["css"] = "dist/styles.css"
		StaticAssets["js"] = "dist/main.js"

		if *localImgServer {
			imgStoragePath = os.Getenv("DEV_IMG_DISK_STORAGE")
			fileStorage = &img.FileStorage{RootPath: imgStoragePath}
		} else {
			client := S3Client(
				os.Getenv("DEV_S3_ENDPOINT"),
				os.Getenv("DEV_S3_KEY"),
				os.Getenv("DEV_S3_SECRET"),
			)
			fileStorage = &img.S3Storage{
				Client:     client,
				BucketName: os.Getenv("DEV_S3_BUCKET"),
			}
		}
	} else {
		infoLog.Println("Production env selected")
		dbUrl = os.Getenv("DB_URL")
		imgBaseUrl = os.Getenv("IMG_URL")
		origin = os.Getenv("ORIGIN")
		client := S3Client(
			os.Getenv("S3_ENDPOINT"),
			os.Getenv("S3_KEY"),
			os.Getenv("S3_SECRET"),
		)

		err = loadAssets()
		if err != nil {
			log.Fatal("Error loading manifest found in production build")
		}
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
	sessionManager.Lifetime = 90 * 24 * time.Hour
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
		moments:        &models.MomentModel{DB: dbpool},
		people:         &models.PersonModel{DB: dbpool},
		profile:        &models.ProfileModel{DB: dbpool},
		recurring:      &models.RecurringModel{DB: dbpool},
		shows:          &models.ShowModel{DB: dbpool},
		tags:           &models.TagModel{DB: dbpool},
		users:          &models.UserModel{DB: dbpool},
		sketches:       &models.SketchModel{DB: dbpool},
		series:         &models.SeriesModel{DB: dbpool},
		sessionManager: sessionManager,
		debugMode:      *debug,
		baseImgUrl:     imgBaseUrl,
		assets:         StaticAssets,
		settings: settings{
			pageSize:          24,
			maxSearchResults:  12,
			localImageServer:  *localImgServer,
			localImageStorage: *localImgStorage,
			origin:            origin,
		},
	}
	app.infoLog.Println("ORIGIN: ", origin)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes("./ui/static/", imgStoragePath, *serveStatic),
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
