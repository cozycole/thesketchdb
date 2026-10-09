package users

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	Insert(user *models.User) error
}

type TokenRepository interface {
	DeleteAllForUser(scope string, userID int) error
	New(userID int, ttl time.Duration, scope string) (*models.Token, error)
}

type Mailer interface {
	Send(to string, template string, data any) error
}

type Config struct {
	origin             string
	activationTokenTTL time.Duration
	passwordTokenTTL   time.Duration
}

type UserService struct {
	users  UserRepository
	tokens TokenRepository
	mailer Mailer
	config Config
}

func New(
	users UserRepository,
	tokens TokenRepository,
	mailer Mailer,
	config Config,
) *UserService {
	return &UserService{
		users,
		tokens,
		mailer,
		config,
	}
}
