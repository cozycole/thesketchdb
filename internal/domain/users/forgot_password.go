package users

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (s *UserService) ForgotPassword(email string) error {
	user, err := s.users.GetByEmail(email)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// we simply don't send an email if no account exists
			return nil
		}
		return err
	}

	token, err := s.tokens.New(
		*user.ID,
		s.config.PasswordTokenTTL,
		models.ScopePasswordReset,
	)
	if err != nil {
		return err
	}

	resetUrl := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.config.Origin,
		*token.Plaintext,
	)

	return s.mailer.Send(
		email,
		"forgot_password.gohtml",
		map[string]any{
			"ResetURL": resetUrl,
		},
	)
}

func (s *UserService) ResetPassword(token, password string) error {
	user, err := s.users.GetByToken(token, models.ScopePasswordReset)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return ErrInvalidToken
		}
		return err
	}

	user.Activated = utils.Ptr(true)
	user.Password.Set(password)
	err = s.users.Update(user)
	if err != nil {
		return err
	}

	// delete token
	err = s.tokens.DeleteAllForUser(models.ScopePasswordReset, *user.ID)
	return err
}
