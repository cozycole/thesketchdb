package users

import (
	"errors"

	"sketchdb.cozycole.net/internal/models"
)

type SocialAuthStatus int

const (
	SocialAuthAuthenticated SocialAuthStatus = iota
	SocialAuthRegistrationRequired
)

type SocialAuthResult struct {
	Status SocialAuthStatus
	User   *models.User
}

func (s *UserService) AuthenticateSocial(provider, providerUserId, email string) (*SocialAuthResult, error) {
	// 1. Check if Google identity already exists.
	user, err := s.users.GetByIdentity("google", providerUserId)

	if err == nil {
		return &SocialAuthResult{
			Status: SocialAuthAuthenticated,
			User:   user,
		}, nil
	}

	if !errors.Is(err, models.ErrNoRecord) {
		return nil, err
	}

	// 2. Check if a user exists with the same email.
	user, err = s.users.GetByEmail(email)

	if err == nil {
		// Automatically link the Google identity.
		err = s.users.CreateIdentity(*user.ID, "google", providerUserId)
		if err != nil {
			return nil, err
		}

		return &SocialAuthResult{
			Status: SocialAuthAuthenticated,
			User:   user,
		}, nil
	}

	if !errors.Is(err, models.ErrNoRecord) {
		return nil, err
	}

	// 3. No matching user or identity.
	return &SocialAuthResult{
		Status: SocialAuthRegistrationRequired,
	}, nil
}

type SocialSignupInput struct {
	Username     string
	Email        string
	ProviderId   string
	ProviderName string
}

// func (s *UserService) CompleteSocialSignup()
