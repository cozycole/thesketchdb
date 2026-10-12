package users

import (
	"errors"

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
