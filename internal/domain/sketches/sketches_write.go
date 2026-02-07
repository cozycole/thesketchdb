package sketches

import (
	"errors"
	"io"
	"mime/multipart"

	"sketchdb.cozycole.net/internal/media"
	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (s *SketchService) CreateSketch(sketch *models.Sketch, thumbnail *multipart.FileHeader) (*models.Sketch, error) {
	if sketch.Episode != nil && sketch.Episode.ID != nil {
		exists, err := s.Repos.Shows.EpisodeExists(safeDeref(sketch.Episode.ID))
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, models.ErrNoEpisode
		}
	}

	if sketch.Creator != nil {
		exists, err := s.Repos.Creators.Exists(*sketch.Creator.ID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, models.ErrNoCreator
		}
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

	err = media.RunImagePipeline(
		thumbnailFile,
		media.Large,
		media.Thumbnail,
		thumbName,
		"/sketch",
		s.ImgStore,
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

func (s *SketchService) UpdateSketch(sketch *models.Sketch, thumbnail []byte) (*models.Sketch, error) {
	oldSketch, err := s.Repos.Sketches.GetById(safeDeref(sketch.ID))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return sketch, models.ErrNoSketch
		} else {
			return sketch, err
		}
	}

	if sketch.Episode != nil && sketch.Episode.ID != nil {
		exists, err := s.Repos.Shows.EpisodeExists(*sketch.Episode.ID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, models.ErrNoEpisode
		}

	}

	if sketch.Creator != nil && sketch.Creator.ID != nil {
		exists, err := s.Repos.Creators.Exists(*sketch.Creator.ID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, models.ErrNoCreator
		}
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

	if thumbnail != nil && oldSketch.ThumbnailName != nil {
		err = media.DeleteImageVariants(s.ImgStore, "sketch", *oldSketch.ThumbnailName)
		if err != nil {
			return sketch, err
		}
	}

	return s.Repos.Sketches.GetById(*sketch.ID)
}
