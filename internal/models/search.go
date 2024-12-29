package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Logic for comprehensive search of all resources

type SearchResult struct {
	Type        *string
	ID          *int
	Name        *string
	Slug        *string
	Img         *string
	UploadDate  *time.Time
	Creator     *string
	CreatorSlug *string
	CreatorImg  *string
	Rank        *float32
}

type SearchModel struct {
	DB *pgxpool.Pool
}

type SearchModelInterface interface {
	Search(query string) ([]*SearchResult, error)
}

func (m *SearchModel) Search(query string) ([]*SearchResult, error) {
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

	UNION ALL

	SELECT 'video' AS type, 
		v.id, 
		v.title AS name, 
		v.slug, 
		v.thumbnail_name AS img, 
		v.upload_date, 
		c.name AS creator, 
		c.slug AS creator_slug, 
		c.profile_img AS creator_img,
		ts_rank(v.search_vector, plainto_tsquery('english', $1)) AS rank
	FROM video as v
	LEFT JOIN video_creator_rel as vcr
	ON v.id = vcr.video_id
	LEFT JOIN creator as c
	ON vcr.creator_id = c.id
	WHERE v.search_vector @@ plainto_tsquery('english', $1)
	ORDER BY rank DESC, name ASC;
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

	results := []*SearchResult{}
	for rows.Next() {
		sr := &SearchResult{}
		err := rows.Scan(
			&sr.Type, &sr.ID, &sr.Name, &sr.Slug, &sr.Img, &sr.UploadDate,
			&sr.Creator, &sr.CreatorSlug, &sr.CreatorImg, &sr.Rank,
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
