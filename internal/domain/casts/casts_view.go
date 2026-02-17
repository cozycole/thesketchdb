package casts

import (
	// "errors"

	"sketchdb.cozycole.net/internal/models"
)

type AdminCastView struct {
	Cast        []*models.CastMember
	Screenshots []*models.CastScreenshot
}

func (s *CastService) GetAdminCast(sketchId int) (AdminCastView, error) {
	cv := AdminCastView{}
	cast, err := s.Repos.Cast.GetCastMembers(sketchId)
	if err != nil {
		return cv, err
	}

	screenshots, err := s.Repos.Cast.GetCastScreenshots(sketchId)
	if err != nil {
		return cv, err
	}

	cv.Cast = cast
	cv.Screenshots = screenshots

	return cv, nil
}
