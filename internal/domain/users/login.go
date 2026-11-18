package users

import (
	"errors"

	"github.com/cozycole/thesketchdb/internal/models"
	"github.com/cozycole/thesketchdb/internal/utils"
)

type LoginInput struct {
	Email    string
	Password string
}

func (s *UserService) Login(input LoginInput) (*models.User, error) {
	user, err := s.users.Authenticate(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return nil, ErrInvalidLoginCredentials
		}
		return nil, err
	}

	if !utils.SafeDeref(user.Activated) {
		return nil, ErrUserNotActivated
	}
	return user, nil
}
