package sketches

import (
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
		sketch.Creator, _ = s.Repos.Creators.GetById(safeDeref(sketch.Creator.ID))
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
		err = s.Repos.Sketches.InsertSketchCreatorRelation(id, *sketch.Creator.ID)
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
