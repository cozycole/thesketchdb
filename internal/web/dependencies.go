package web

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/mail"
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
