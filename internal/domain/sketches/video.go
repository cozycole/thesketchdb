package sketches

import "sketchdb.cozycole.net/internal/models"

func (s *SketchService) AddVideo(sketchId int, hotS3Key string) (*models.SketchVideo, error) {
	video := &models.SketchVideo{HotS3Key: &hotS3Key}

	err := s.Repos.Sketches.InsertSketchVideo(sketchId, video)
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (s *SketchService) GetVideos(sketchId int) ([]*models.SketchVideo, error) {
	videos, err := s.Repos.Sketches.GetVideos(sketchId)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
