package casts

import (
	"fmt"

	"sketchdb.cozycole.net/internal/media"
	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (s *CastService) ReorderCast(sketchId int, castIds []int) error {
	cast, err := s.Repos.Cast.GetCastMembers(sketchId)
	if err != nil {
		return err
	}

	err = validateCastIds(castIds, cast)
	if err != nil {
		return err
	}

	return s.Repos.Cast.UpdatePositions(castIds)
}

func (s *CastService) CreateCastMember(cm *models.CastMember, thumbnail []byte, profile []byte) (*models.CastMember, error) {
	if cm.SketchID == nil {
		return nil, fmt.Errorf("sketch id not defined in cast member input")
	}

	if thumbnail != nil {
		thumbName, err := media.GenerateFileName(thumbnail)
		if err != nil {
			return nil, err
		}
		cm.ThumbnailName = &thumbName
	}

	if profile != nil {
		profileName, err := media.GenerateFileName(profile)
		if err != nil {
			return nil, err
		}
		cm.ProfileImg = &profileName
	}

	// determine postion order for display
	cast, err := s.Repos.Cast.GetCastMembers(*cm.SketchID)
	if err != nil {
		return nil, err
	}

	position := getNextPosition(cast)
	cm.Position = &position

	err = s.Repos.Cast.Insert(*cm.SketchID, cm)
	if err != nil {
		return nil, err
	}

	if thumbnail != nil {
		err = media.RunImagePipeline(
			thumbnail,
			media.Medium,
			media.Thumbnail,
			*cm.ThumbnailName,
			"/cast/thumbnail",
			s.ImgStore,
		)

		if err != nil {
			s.Repos.Cast.Delete(*cm.ID)
			return nil, err
		}
	}

	if profile != nil {
		err = media.RunImagePipeline(
			profile,
			media.Medium,
			media.Profile,
			*cm.ProfileImg,
			"/cast/profile",
			s.ImgStore,
		)

		if err != nil {
			s.Repos.Cast.Delete(*cm.ID)
			return nil, err
		}
	}

	newMember, err := s.Repos.Cast.GetById(*cm.ID)
	if err != nil {
		return nil, err
	}

	return newMember, nil
}

func (s *CastService) UpdateCastMember(cm *models.CastMember, thumbnail []byte, profile []byte) (*models.CastMember, error) {
	if cm.ID == nil {
		return nil, fmt.Errorf("no id specified for cast member update")
	}

	// get stale cast member for file names
	staleCast, err := s.Repos.Cast.GetById(*cm.ID)
	if err != nil {
		return nil, err
	}

	currentThumbnailName := utils.SafeDeref(staleCast.ThumbnailName)
	currentProfileName := utils.SafeDeref(staleCast.ProfileImg)
	if thumbnail != nil {
		newThumbnailName, err := media.GenerateFileName(thumbnail)
		if err != nil {
			return nil, err
		}

		err = media.RunImagePipeline(
			thumbnail,
			media.Medium,
			media.Thumbnail,
			newThumbnailName,
			"/cast/thumbnail",
			s.ImgStore,
		)

		currentThumbnailName = newThumbnailName
	}
	if profile != nil {
		newProfileName, err := media.GenerateFileName(profile)
		if err != nil {
			return nil, err
		}

		err = media.RunImagePipeline(
			profile,
			media.Medium,
			media.Profile,
			newProfileName,
			"/cast/profile",
			s.ImgStore,
		)

		currentProfileName = newProfileName
	}

	cm.ThumbnailName = &currentThumbnailName
	cm.ProfileImg = &currentProfileName

	err = s.Repos.Cast.Update(cm)
	if err != nil {
		return nil, err
	}

	if thumbnail != nil && staleCast.ThumbnailName != nil {
		media.DeleteImageVariants(s.ImgStore, "cast/thumbnail", *staleCast.ThumbnailName)
	}
	if profile != nil && staleCast.ProfileImg != nil {
		media.DeleteImageVariants(s.ImgStore, "cast/profile", *staleCast.ProfileImg)
	}

	newMember, err := s.Repos.Cast.GetById(*cm.ID)
	if err != nil {
		return nil, err
	}

	return newMember, nil
}

func (s *CastService) DeleteCastmember(id int) error {
	castMember, err := s.Repos.Cast.GetById(id)
	if err != nil {
		return err
	}

	err = s.Repos.Cast.Delete(id)
	if err != nil {
		return err
	}

	if utils.SafeDeref(castMember.ThumbnailName) != "" {
		media.DeleteImageVariants(s.ImgStore, "cast/thumbnail", *castMember.ThumbnailName)
	}
	if utils.SafeDeref(castMember.ProfileImg) != "" {
		media.DeleteImageVariants(s.ImgStore, "cast/profile", *castMember.ProfileImg)
	}

	return nil
}
