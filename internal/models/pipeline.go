package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PipelineJob struct {
	ID     *int    `json:"id"`
	Status *string `json:"status"`
	Error  *string `json:"error"`
}

type PipelineModelInterface interface {
	Insert(int, *PipelineJob) error
}

type PipelineModel struct {
	DB *pgxpool.Pool
}

func (m *PipelineModel) Insert(videoId int, pipeline *PipelineJob) error {
	stmt := `
		INSERT INTO pipeline_jobs (video_id, status, error)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	result := m.DB.QueryRow(
		context.Background(), stmt, videoId, pipeline.Status, pipeline.Error,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return err
	}
	pipeline.ID = &id

	return nil
}
