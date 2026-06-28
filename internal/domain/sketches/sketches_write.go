package sketches

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"sketchdb.cozycole.net/internal/media"
	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (s *SketchService) CreateSketch(sketch *models.Sketch, thumbnail *multipart.FileHeader, cropBorder bool) (*models.Sketch, error) {
	if sketch.Episode != nil && sketch.Episode.ID != nil {
		ep, err := s.Repos.Shows.GetEpisode(*sketch.Episode.ID)
		if err != nil {
			return nil, err
		}

		sketch.Episode = ep
	}

	if sketch.Creator != nil && sketch.Creator.ID != nil {
		c, err := s.Repos.Creators.GetById(*sketch.Creator.ID)
		if err != nil {
			return nil, err
		}

		sketch.Creator = c.ToRef()
	}

	slug := createSketchSlug(sketch)
	sketch.Slug = &slug

	thumbName, err := utils.GenerateFileName(thumbnail)
	if err != nil {
		return nil, err
	}
	sketch.ThumbnailName = &thumbName

	youtubeID, _ := extractYouTubeVideoID(*sketch.URL)
	if youtubeID != "" {
		sketch.YoutubeID = &youtubeID
	}

	f, err := thumbnail.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// extract bytes from multipart.File
	thumbnailFile, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	id, err := s.Repos.Sketches.Insert(sketch)
	if err != nil {
		return nil, err
	}

	if sketch.Creator != nil {
		err = s.Repos.Sketches.SyncSketchCreators(*sketch.ID, []int{safeDeref(sketch.Creator.ID)})
		if err != nil {
			s.Repos.Sketches.Delete(id)
			return nil, err
		}
	}

	if len(sketch.Tags) > 0 {
		err = s.Repos.Sketches.BatchUpdateTags(id, sketch.Tags)
		if err != nil {
			return nil, err
		}
	}

	err = media.RunImagePipeline(
		thumbnailFile,
		media.Large,
		media.Thumbnail,
		thumbName,
		"/sketch",
		s.ImgStore,
		cropBorder,
	)
	if err != nil {
		s.Repos.Sketches.Delete(id)
		return nil, err
	}

	createdSketch, err := s.Repos.Sketches.GetById(id)
	if err != nil {
		s.Repos.Sketches.Delete(id)
		return nil, err
	}

	return createdSketch, nil
}

func (s *SketchService) UpdateSketch(sketch *models.Sketch, thumbnail []byte, cropBorder bool) (*models.Sketch, error) {
	oldSketch, err := s.Repos.Sketches.GetById(safeDeref(sketch.ID))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return sketch, models.ErrNoSketch
		} else {
			return sketch, err
		}
	}

	if sketch.Episode != nil && sketch.Episode.ID != nil {
		ep, err := s.Repos.Shows.GetEpisode(*sketch.Episode.ID)
		if err != nil {
			return nil, err
		}

		sketch.Episode = ep
	}

	if sketch.Creator != nil && sketch.Creator.ID != nil {
		c, err := s.Repos.Creators.GetById(*sketch.Creator.ID)
		if err != nil {
			return nil, err
		}

		sketch.Creator = c.ToRef()
	}

	slug := createSketchSlug(sketch)
	sketch.Slug = &slug

	youtubeID, _ := extractYouTubeVideoID(*sketch.URL)
	if youtubeID != "" {
		sketch.YoutubeID = &youtubeID
	}

	// try to save new file first
	// if theres an error, we don't want the sketch to be updated
	// with a new thumbnail name that now doesn't exist
	thumbnailName := safeDeref(oldSketch.ThumbnailName)
	if thumbnail != nil {
		var err error
		thumbnailName, err = media.GenerateFileName(thumbnail)
		if err != nil {
			return sketch, err
		}

		err = media.RunImagePipeline(
			thumbnail,
			media.Large,
			media.Thumbnail,
			thumbnailName,
			"/sketch",
			s.ImgStore,
			cropBorder,
		)
		if err != nil {
			return sketch, err
		}
	}

	sketch.ThumbnailName = &thumbnailName
	err = s.Repos.Sketches.Update(sketch)
	if err != nil {
		return sketch, err
	}

	if sketch.Creator != nil {
		err = s.Repos.Sketches.SyncSketchCreators(*sketch.ID, []int{safeDeref(sketch.Creator.ID)})
		if err != nil {
			return nil, err
		}
	}

	err = s.Repos.Sketches.BatchUpdateTags(*sketch.ID, sketch.Tags)
	if err != nil {
		return nil, err
	}

	if thumbnail != nil && oldSketch.ThumbnailName != nil {
		err = media.DeleteImageVariants(s.ImgStore, "sketch", *oldSketch.ThumbnailName)
		if err != nil {
			return sketch, err
		}
	}

	return s.GetSketch(*sketch.ID)
}

type DeleteSketchInfo struct {
	Videos          []*models.SketchVideo
	CastScreenshots []*models.CastScreenshot
}

func (s *SketchService) DeleteSketch(id int) (*DeleteSketchInfo, error) {
	deleteInfo := DeleteSketchInfo{}

	var err error
	deleteInfo.Videos, err = s.Repos.Sketches.GetVideos(id)
	if err != nil {
		return nil, err
	}

	deleteInfo.CastScreenshots, err = s.Repos.Cast.GetCastScreenshots(id)
	if err != nil {
		return nil, err
	}

	err = s.Repos.Sketches.Delete(id)
	if err != nil {
		return nil, err
	}

	return &deleteInfo, nil
}

// Delete any images or videos associated with a previously deleted sketch
func (s *SketchService) CleanupSketchMedia(info *DeleteSketchInfo) error {
	coldVidKeys := []string{}
	hotVidKeys := []string{}
	for _, v := range info.Videos {
		if v.HotS3Key != nil {
			hotKey := fmt.Sprintf("video/%s", *v.HotS3Key)
			hotVidKeys = append(hotVidKeys, hotKey)
		}
		if v.ColdS3Key != nil {
			coldKey := fmt.Sprintf("video/%s", *v.ColdS3Key)
			coldVidKeys = append(coldVidKeys, coldKey)
		}
	}

	screenshotKeys := []string{}
	for _, s := range info.CastScreenshots {
		if s.ProfileImage != nil {
			shotKey := fmt.Sprintf("cast_auto_screenshots/profile/%s", *s.ProfileImage)
			screenshotKeys = append(screenshotKeys, shotKey)
		}
		if s.ThumbnailName != nil {
			shotKey := fmt.Sprintf("cast_auto_screenshots/thumbnail/%s", *s.ThumbnailName)
			screenshotKeys = append(screenshotKeys, shotKey)
		}
	}

	err := s.ArchiveStore.DeleteFiles(coldVidKeys)
	if err != nil {
		return err
	}
	err = s.ImgStore.DeleteFiles(hotVidKeys)
	if err != nil {
		return err
	}
	err = s.ImgStore.DeleteFiles(screenshotKeys)
	return err
}

func (s *SketchService) DeleteScreenshots(sketchId int) error {
	screenshots, err := s.Repos.Cast.GetCastScreenshots(sketchId)
	if err != nil {
		return err
	}

	err = s.Repos.Sketches.DeleteScreenshots(sketchId)
	if err != nil {
		return err
	}

	screenshotKeys := []string{}
	for _, s := range screenshots {
		if s.ProfileImage != nil {
			shotKey := fmt.Sprintf("cast_auto_screenshots/profile/%s", *s.ProfileImage)
			screenshotKeys = append(screenshotKeys, shotKey)
		}
		if s.ThumbnailName != nil {
			shotKey := fmt.Sprintf("cast_auto_screenshots/thumbnail/%s", *s.ThumbnailName)
			screenshotKeys = append(screenshotKeys, shotKey)
		}
	}

	return s.ImgStore.DeleteFiles(screenshotKeys)
}
