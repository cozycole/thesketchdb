package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A profile a generic struct that represents
// either a person, character, creator or user

type ProfileResult struct {
	Type *string
	ID   *int
	Name *string
	Slug *string
	Img  *string
	Date *time.Time
	Rank *float32
}

type ProfileModel struct {
	DB *pgxpool.Pool
}

type ProfileModelInterface interface {
	Search(query string) ([]*ProfileResult, error)
}

func (m *ProfileModel) Search(query string) ([]*ProfileResult, error) {
	stmt := `
	SELECT 'person' AS type, 
       id, 
       CONCAT(first, ' ', last) AS name, 
       slug, 
       profile_img AS img, 
       CAST(NULL as date) AS upload_date, 
       NULL AS creator, 
       NULL AS creator_slug, 
	   NULL AS creator_img,
       ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
	FROM person
	WHERE search_vector @@ plainto_tsquery('english', $1)

	UNION ALL

	SELECT 'character' AS type, 
		id, 
		name, 
		slug, 
		img_name AS img, 
		NULL AS upload_date, 
		NULL AS creator, 
		NULL AS creator_slug, 
		NULL AS creator_img,
		ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
	FROM character
	WHERE search_vector @@ plainto_tsquery('english', $1)

	UNION ALL

	SELECT 'creator' AS type, 
		id, 
		name, 
		slug, 
		profile_img AS img, 
		NULL AS upload_date, 
		NULL AS creator, 
		NULL AS creator_slug, 
		NULL AS creator_img,
		ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
	FROM creator
	WHERE search_vector @@ plainto_tsquery('english', $1)
	`

	rows, err := m.DB.Query(context.Background(), stmt, query)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	results := []*ProfileResult{}
	for rows.Next() {
		sr := &ProfileResult{}
		err := rows.Scan(
			&sr.Type, &sr.ID, &sr.Name, &sr.Slug, &sr.Img, &sr.Rank,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, sr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
