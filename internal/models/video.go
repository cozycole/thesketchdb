package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Video struct {
	ID         int
	Title      string
	URL        string
	Thumbnail  string
	Rating     string
	UploadDate time.Time
	Creator    *Creator
}

type VideoModelInterface interface {
	// Insert(title string, url string, thumbnail string, rating, uploadDate time.Time)
	Search(search string, offset int) ([]*Video, error)
	GetAll(offset int) ([]*Video, error)
}

type VideoModel struct {
	DB         *pgxpool.Pool
	ResultSize int
}

func (m *VideoModel) Search(search string, offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.thumbnail_name, v.creation_date,
		c.id, c.name, c.profile_img_path
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE search_vector @@ to_tsquery('english', $1)
		LIMIT $2
		OFFSET $3;
	`
	rows, err := m.DB.Query(context.Background(), stmt, search, m.ResultSize, m.ResultSize*offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		err := rows.Scan(&v.ID, &v.Title, &v.URL, &v.Thumbnail)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

// Will make DRY later
func (m *VideoModel) GetAll(offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.thumbnail_name, v.creation_date,
		c.id, c.name, c.profile_img_path
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		LIMIT $1
		OFFSET $2;
	`
	rows, err := m.DB.Query(context.Background(), stmt, m.ResultSize, m.ResultSize*offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		c := &Creator{}
		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Thumbnail, &v.UploadDate,
			&c.ID, &c.Name, &c.ProfileImage,
		)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		videos = append(videos, v)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}
