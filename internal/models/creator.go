package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Creator struct {
	ID              *int
	Name            *string
	URL             *string
	ProfileImage    *string
	Slug            *string
	EstablishedDate *time.Time
}

type CreatorModelInterface interface {
	Insert(name, url, imgName, imgExt string, establishedDate time.Time) (int, string, string, error)
	Get(id int) (*Creator, error)
	Exists(id int) (bool, error)
	GetBySlug(slug string) (*Creator, error)
	Search(query string) ([]*Creator, error)
	SearchCount(query string) (int, error)
	VectorSearch(query string) ([]*ProfileResult, error)
}

type CreatorModel struct {
	DB *pgxpool.Pool
}

func (m *CreatorModel) GetBySlug(slug string) (*Creator, error) {
	stmt := `SELECT id, name, slug, page_url, profile_img, date_established FROM creator
	WHERE slug = $1`

	row := m.DB.QueryRow(context.Background(), stmt, slug)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.URL, &c.ProfileImage, &c.EstablishedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return c, nil
}

func (m *CreatorModel) Insert(name, url, slug, imgExt string, establishedDate time.Time) (int, string, string, error) {
	stmt := `
	INSERT INTO creator (name, page_url, date_established, slug, profile_img)
	VALUES ($1,$2,$3,
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('creator', 'id'))),
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('creator', 'id')), $5::text))
	RETURNING id, slug, profile_img;`

	var id int
	var fullImgName string

	row := m.DB.QueryRow(context.Background(), stmt, name, url, establishedDate, slug, imgExt)
	err := row.Scan(&id, &slug, &fullImgName)
	if err != nil {
		return 0, "", "", err
	}
	return id, slug, fullImgName, err
}

func (m *CreatorModel) Search(query string) ([]*Creator, error) {
	query = query + "%"
	stmt := `SELECT c.id, c.slug, c.name, c.profile_img
			FROM creator as c
			WHERE name ILIKE $1
			ORDER BY name`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	creators := []*Creator{}
	for rows.Next() {
		c := &Creator{}
		err := rows.Scan(
			&c.ID, &c.Slug, &c.Name, &c.ProfileImage,
		)
		if err != nil {
			return nil, err
		}
		creators = append(creators, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return creators, nil
}

func (m *CreatorModel) Get(id int) (*Creator, error) {
	stmt := `SELECT id, name, slug, page_url, profile_img, date_established FROM creator
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.URL, &c.ProfileImage, &c.EstablishedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil
}

func (m *CreatorModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM creator WHERE id = $1`
	row := m.DB.QueryRow(context.Background(), stmt, id)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (m *CreatorModel) VectorSearch(query string) ([]*ProfileResult, error) {
	stmt := `
		SELECT id, name, profile_img, slug, date_established, ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
		FROM creator
		WHERE search_vector @@ plainto_tsquery('english', $1)
		ORDER BY rank desc
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
		pr := &ProfileResult{}
		err := rows.Scan(
			&pr.ID, &pr.Name, &pr.Img, &pr.Slug, &pr.Date, &pr.Rank,
		)
		if err != nil {
			return nil, err
		}

		resType := "creator"
		pr.Type = &resType

		results = append(results, pr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *CreatorModel) SearchCount(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM creator
		WHERE search_vector @@ plainto_tsquery('english', $1)
	`
	var count int
	row := m.DB.QueryRow(context.Background(), stmt, query)
	err := row.Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	return count, nil
}
