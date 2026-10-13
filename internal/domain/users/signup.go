package users

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

type SignupInput struct {
	Username string
	Email    string
	Password string
}

func (s *UserService) Signup(input SignupInput) error {
	existingUser, err := s.users.GetByEmail(input.Email)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}

	// if an activated user with this email already exists
	// we do nothing and simply want to redirect them to the
	// "Check your email" page
	if existingUser != nil && utils.SafeDeref(existingUser.Activated) {
		return nil
	}

	// we insert the user if no user with the provided email exists
	var newUser models.User
	if existingUser == nil {
		newUser = models.User{
			Username: &input.Username,
			Email:    &input.Email,
		}

		err = newUser.Password.Set(input.Password)
		if err != nil {
			return err
		}

		err = s.users.Insert(&newUser)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateUsername) {
				return ErrUsernameUnavailable
			}
			return err
		}
	}

	// if an inactivated user exists already, we delete any activation tokens for them
	userId := utils.SafeDeref(newUser.ID)
	userEmail := utils.SafeDeref(newUser.Email)
	if existingUser != nil {
		err = s.tokens.DeleteAllForUser(models.ScopeActivation, *existingUser.ID)
		if err != nil {
			return err
		}
		userId = *existingUser.ID
		userEmail = *existingUser.Email
	}

	// now create a token either for the previously inactivated user or
	// the newly inserted user
	token, err := s.tokens.New(userId, s.config.ActivationTokenTTL, models.ScopeActivation)
	if err != nil {
		return err
	}

	return s.mailer.Send(
		userEmail,
		"user_welcome.gohtml",
		map[string]string{
			"ActivationURL": fmt.Sprintf(
				"%s/verify-email?token=%s",
				s.config.Origin,
				*token.Plaintext,
			),
		},
	)
}
