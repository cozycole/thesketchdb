package main

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"sketchdb.cozycole.net/internal/domain/casts"
	"sketchdb.cozycole.net/internal/domain/characters"
	"sketchdb.cozycole.net/internal/domain/creators"
	"sketchdb.cozycole.net/internal/domain/people"
	"sketchdb.cozycole.net/internal/domain/pipeline"
	"sketchdb.cozycole.net/internal/domain/quotes"
	"sketchdb.cozycole.net/internal/domain/recurring"
	"sketchdb.cozycole.net/internal/domain/series"
	"sketchdb.cozycole.net/internal/domain/shows"
	"sketchdb.cozycole.net/internal/domain/sketches"
	"sketchdb.cozycole.net/internal/domain/tags"
	"sketchdb.cozycole.net/internal/domain/users"

	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/mail"
	"sketchdb.cozycole.net/internal/models"
)

type Services struct {
	Casts      casts.CastService
	Characters characters.CharacterService
	Creators   creators.CreatorService
	People     people.PersonService
	Pipeline   pipeline.PipelineService
	Quotes     quotes.QuoteService
	Recurring  recurring.RecurringService
	Series     series.SeriesService
	Shows      shows.ShowService
	Sketches   sketches.SketchService
	Tags       tags.TagsService
	Users      users.UserService
}

func NewServices(
	repos models.Repositories,
	fileStore fileStore.FileStorageInterface,
	archiveStore fileStore.FileStorageInterface,
	mailer mail.MailerInterface,
	config AppConfig,
) Services {
	return Services{
		Sketches: sketches.SketchService{
			Repos:        repos,
			ImgStore:     fileStore,
			ArchiveStore: archiveStore,
		},
		Creators: creators.CreatorService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Shows: shows.ShowService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Recurring: recurring.RecurringService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Series: series.SeriesService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Casts: casts.CastService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		People: people.PersonService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Pipeline: pipeline.PipelineService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Characters: characters.CharacterService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Quotes: quotes.QuoteService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Tags: tags.TagsService{
			Repos:    repos,
			ImgStore: fileStore,
		},
		Users: *users.New(repos.Users, repos.Tokens, mailer,
			users.Config{
				Origin:             config.Origin,
				ActivationTokenTTL: config.ActivationTTL,
				PasswordTokenTTL:   config.PasswordTTL,
			}),
	}
}

func newRepositories(dbpool *pgxpool.Pool) models.Repositories {
	return models.Repositories{
		Cast:       &models.CastModel{DB: dbpool},
		Categories: &models.CategoryModel{DB: dbpool},
		Characters: &models.CharacterModel{DB: dbpool},
		Creators:   &models.CreatorModel{DB: dbpool},
		Quotes:     &models.QuoteModel{DB: dbpool},
		People:     &models.PersonModel{DB: dbpool},
		Profile:    &models.ProfileModel{DB: dbpool},
		Pipeline:   &models.PipelineModel{DB: dbpool},
		Recurring:  &models.RecurringModel{DB: dbpool},
		Shows:      &models.ShowModel{DB: dbpool},
		Tags:       &models.TagModel{DB: dbpool},
		Tokens:     &models.TokenModel{DB: dbpool},
		Users:      &models.UserModel{DB: dbpool},
		Sketches:   &models.SketchModel{DB: dbpool},
		Series:     &models.SeriesModel{DB: dbpool},
	}
}
