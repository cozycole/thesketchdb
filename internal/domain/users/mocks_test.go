package users

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

// -------------------------
// User repository
// -------------------------

type spyUserRepository struct {
	// Stub configuration
	foundUser *models.User
	findErr   error
	insertErr error
	nextID    int

	// Recorded interactions
	getByEmailCalls []string
	insertedUsers   []*models.User
}

func (s *spyUserRepository) GetByEmail(email string) (*models.User, error) {
	s.getByEmailCalls = append(s.getByEmailCalls, email)

	return s.foundUser, s.findErr
}

func (s *spyUserRepository) Insert(user *models.User) error {
	s.insertedUsers = append(s.insertedUsers, user)

	if s.insertErr != nil {
		return s.insertErr
	}

	// Simulate the database assigning an ID.
	user.ID = &s.nextID

	return nil
}

// -------------------------
// Token repository
// -------------------------

type deleteTokenCall struct {
	scope  string
	userID int
}

type createTokenCall struct {
	userID int
	ttl    time.Duration
	scope  string
}

type spyTokenRepository struct {
	// Stub configuration
	deleteErr error
	createErr error

	// Recorded interactions
	deleteCalls []deleteTokenCall
	createCalls []createTokenCall
}

func (s *spyTokenRepository) DeleteAllForUser(
	scope string,
	userID int,
) error {
	s.deleteCalls = append(s.deleteCalls, deleteTokenCall{
		scope:  scope,
		userID: userID,
	})

	return s.deleteErr
}

func (s *spyTokenRepository) New(
	userID int,
	ttl time.Duration,
	scope string,
) (*models.Token, error) {
	s.createCalls = append(s.createCalls, createTokenCall{
		userID: userID,
		ttl:    ttl,
		scope:  scope,
	})

	if s.createErr != nil {
		return nil, s.createErr
	}

	plaintext := "test-token"

	return &models.Token{
		Plaintext: &plaintext,
	}, nil
}

// -------------------------
// Mailer
// -------------------------

type mailCall struct {
	to       string
	template string
	data     any
}

type spyMailer struct {
	// Stub configuration
	sendErr error

	// Recorded interactions
	calls []mailCall
}

func (s *spyMailer) Send(
	to string,
	template string,
	data any,
) error {
	s.calls = append(s.calls, mailCall{
		to:       to,
		template: template,
		data:     data,
	})

	return s.sendErr
}
