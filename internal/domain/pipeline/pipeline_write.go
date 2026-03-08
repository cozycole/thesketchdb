package pipeline

import "sketchdb.cozycole.net/internal/models"

func (s *PipelineService) AddPipelineJob(videoId int) (*models.PipelineJob, error) {
	status := "pending"
	job := &models.PipelineJob{
		Status: &status,
	}

	err := s.Repos.Pipeline.Insert(videoId, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}
