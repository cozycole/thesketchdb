package casts

import (
	// "errors"
	"fmt"

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

type CastsListResult struct {
	Casts         []*models.CastMember
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	PersonRefs    []*models.PersonRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      models.Metadata
	Filter        *models.Filter
}

func (s *CastService) ListCasts(f *models.Filter, includeRefs bool) (CastsListResult, error) {
	result := CastsListResult{}
	casts, metadata, err := s.Repos.Cast.List(f)
	if err != nil {
		return result, fmt.Errorf("list cast error: %w", err)
	}

	result.Metadata = metadata
	result.Filter = f
	result.Casts = casts
	return result, nil
}
