package main

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/mail"
	"sketchdb.cozycole.net/internal/models"
)

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

func NewApplication(
	config AppConfig,
	deps AppDependencies,
) *Application {
	db := deps.DB

	return &Application{
		Config:         config,
		InfoLog:        deps.InfoLog,
		ErrorLog:       deps.ErrorLog,
		TemplateCache:  deps.TemplateCache,
		SessionManager: deps.SessionManager,
		FormDecoder:    deps.FormDecoder,

		FileStorage:    deps.FileStorage,
		ArchiveStorage: deps.ArchiveStorage,
		Mail:           deps.Mailer,

		Categories: &models.CategoryModel{DB: db},
		Characters: &models.CharacterModel{DB: db},
		Creators:   &models.CreatorModel{DB: db},
		People:     &models.PersonModel{DB: db},
		Profile:    &models.ProfileModel{DB: db},
		Quotes:     &models.QuoteModel{DB: db},
		Recurring:  &models.RecurringModel{DB: db},
		Series:     &models.SeriesModel{DB: db},
		Shows:      &models.ShowModel{DB: db},
		Sketches:   &models.SketchModel{DB: db},
		Tags:       &models.TagModel{DB: db},
		Tokens:     &models.TokenModel{DB: db},
		Users:      &models.UserModel{DB: db},

		Services: NewServices(
			newRepositories(db),
			deps.FileStorage,
			deps.ArchiveStorage,
			deps.Mailer,
			config,
		),

		Assets: deps.Assets,
	}
}
