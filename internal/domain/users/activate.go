package users

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (s *UserService) Activate(token string) (*models.User, error) {
	user, err := s.users.GetByToken(token, models.ScopeActivation)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	user.Activated = utils.Ptr(true)
	err = s.users.Update(user)
	if err != nil {
		return nil, err
	}

	err = s.tokens.DeleteAllForUser(models.ScopeActivation, *user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ResendActivationLink(email string) error {
	user, err := s.users.GetByEmail(email)

	// if they don't exist or they're already activated, we aren't sending an email
	if errors.Is(err, models.ErrNoRecord) || utils.SafeDeref(user.Activated) {
		return nil
	}

	token, err := s.tokens.New(*user.ID, s.config.ActivationTokenTTL, models.ScopeActivation)
	if err != nil {
		return err
	}

	return s.mailer.Send(
		*user.Email,
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
