package main

import (
	"html/template"
	"log"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/mail"
	"sketchdb.cozycole.net/internal/models"
)

type AppDependencies struct {
	DB             *pgxpool.Pool
	templateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
	formDecoder    *form.Decoder
	Mailer         mail.MailerInterface
	FileStorage    fileStore.FileStorageInterface
	ArchiveStorage fileStore.FileStorageInterface
}

type AppConfig struct {
	Origin        string
	DebugMode     bool
	DevEnv        bool
	BaseImgUrl    string
	PasswordTTL   time.Duration
	ActivationTTL time.Duration

	PageSize         int
	MaxSearchResults int
}

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	TemplateCache  map[string]*template.Template
	FileStorage    fileStore.FileStorageInterface
	ArchiveStorage fileStore.FileStorageInterface
	Cast           models.CastModelInterface
	Categories     models.CategoryInterface
	Characters     models.CharacterModelInterface
	Creators       models.CreatorModelInterface
	Quotes         models.QuoteModelInterface
	People         models.PersonModelInterface
	Profile        models.ProfileModelInterface
	Recurring      models.RecurringModelInterface
	Shows          models.ShowModelInterface
	Series         models.SeriesModelInterface
	Tags           models.TagModelInterface
	Tokens         models.TokenModelInterface
	Users          models.UserModelInterface
	Sketches       models.SketchModelInterface
	Services       Services
	Mail           mail.MailerInterface
	SessionManager *scs.SessionManager
	FormDecoder    *form.Decoder
	Assets         map[string]string
	Config         AppConfig
}
