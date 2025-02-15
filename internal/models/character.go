package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Character struct {
	ID          *int
	Slug        *string
	Name        *string
	Image       *string
	Description *string
	Portrayal   *Person
}

type CharacterModelInterface interface {
	Search(search string) ([]*Character, error)
	Exists(id int) (bool, error)
	SearchCount(query string) (int, error)
	VectorSearch(query string, limit, offset int) ([]*ProfileResult, error)
}

type CharacterModel struct {
	DB *pgxpool.Pool
}

func (m *CharacterModel) Search(query string) ([]*Character, error) {
	query = query + "%"
	stmt := `SELECT c.id, c.slug, c.name, c.img_name
			FROM character as c
			WHERE name ILIKE $1
			ORDER BY name`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characters := []*Character{}
	for rows.Next() {
		c := &Character{}
		err := rows.Scan(
			&c.ID, &c.Slug, &c.Name, &c.Image,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return characters, nil
}

func (m *CharacterModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM character WHERE id = $1`
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

func (m *CharacterModel) VectorSearch(query string, limit, offset int) ([]*ProfileResult, error) {
	fmt.Printf("Got here %s %d %d\n", query, limit, offset)
	stmt := `
		SELECT id, name, img_name, slug, ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
		FROM character
		WHERE search_vector @@ plainto_tsquery('english', $1)
		ORDER BY rank desc
		LIMIT $2
		OFFSET $3;
	`

	rows, err := m.DB.Query(context.Background(), stmt, query, limit, offset)

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
			&pr.ID, &pr.Name, &pr.Img, &pr.Slug, &pr.Rank,
		)
		if err != nil {
			return nil, err
		}

		resType := "character"
		pr.Type = &resType

		results = append(results, pr)
	}

	return results, nil
}

func (m *CharacterModel) SearchCount(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM character as c
		WHERE c.search_vector @@ plainto_tsquery('english', $1)
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
