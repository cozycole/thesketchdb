package main

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cozycole/thesketchdb/internal/domain/casts"
	"github.com/cozycole/thesketchdb/internal/domain/characters"
	"github.com/cozycole/thesketchdb/internal/domain/creators"
	"github.com/cozycole/thesketchdb/internal/domain/people"
	"github.com/cozycole/thesketchdb/internal/domain/pipeline"
	"github.com/cozycole/thesketchdb/internal/domain/quotes"
	"github.com/cozycole/thesketchdb/internal/domain/recurring"
	"github.com/cozycole/thesketchdb/internal/domain/series"
	"github.com/cozycole/thesketchdb/internal/domain/shows"
	"github.com/cozycole/thesketchdb/internal/domain/sketches"
	"github.com/cozycole/thesketchdb/internal/domain/tags"
	"github.com/cozycole/thesketchdb/internal/domain/users"

	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/mail"
	"github.com/cozycole/thesketchdb/internal/models"
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
