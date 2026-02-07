package models

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Creator struct {
	ID              *int       `json:"id"`
	Slug            *string    `json:"slug"`
	Name            *string    `json:"name"`
	ProfileImage    *string    `json:"profileImage"`
	Alias           *string    `json:"alias"`
	URL             *string    `json:"url"`
	EstablishedDate *time.Time `json:"establishedDate"`
}

type CreatorRef struct {
	ID           *int    `json:"id"`
	Slug         *string `json:"slug"`
	Name         *string `json:"name"`
	ProfileImage *string `json:"profileImage"`
}

func (c *Creator) HasId() bool {
	return c.ID != nil
}

type CreatorModelInterface interface {
	Delete(id int) error
	Exists(id int) (bool, error)
	GetById(id int) (*Creator, error)
	GetCast(id int) ([]*Person, error)
	GetCount(filter *Filter) (int, error)
	GetCreatorRefs([]int) ([]*CreatorRef, error)
	Insert(creator *Creator) (int, error)
	List(filter *Filter) ([]*CreatorRef, Metadata, error)
	Search(query string) ([]*Creator, error)
	SearchCount(query string) (int, error)
	Update(creator *Creator) error
	VectorSearch(query string) ([]*ProfileResult, error)
}

type CreatorModel struct {
	DB *pgxpool.Pool
}

func (m *CreatorModel) Delete(id int) error {
	stmt := `
		DELETE FROM creator
		WHERE id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *CreatorModel) List(filter *Filter) ([]*CreatorRef, Metadata, error) {
	query := `SELECT count(*) OVER(), c.id, c.name, c.slug, c.profile_img%s
			FROM creator as c
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		rankParam := fmt.Sprintf(`
		, ts_rank(
			setweight(to_tsvector('english', c.name || ' ' || COALESCE(c.alias, '')) , 'A'),
			websearch_to_tsquery('english', $%d)
		) AS rank
		`, argIndex)

		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`AND (
            to_tsvector('english', c.name || ' ' || COALESCE(c.alias, '')) @@ websearch_to_tsquery('english', $%d)
			OR 
			(c.name || ' ' || coalesce(c.alias,'')) ILIKE '%%' || $%d || '%%'
			)
		`, argIndex, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	// fmt.Print(query)
	rows, err := m.DB.Query(context.Background(), query, args...)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, Metadata{}, fmt.Errorf("%s:%d: %w", file, line, err)
	}

	creators := []*CreatorRef{}
	var totalCount int
	for rows.Next() {
		var c CreatorRef
		destinations := []any{
			&totalCount, &c.ID, &c.Name, &c.Slug, &c.ProfileImage,
		}

		var rank *float32
		if filter.Query != "" {
			destinations = append(destinations, &rank)
		}
		err := rows.Scan(destinations...)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			return nil, Metadata{}, fmt.Errorf("%s:%d: %w", file, line, err)
		}

		creators = append(creators, &c)
	}

	if err = rows.Err(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, Metadata{}, fmt.Errorf("%s:%d: %w", file, line, err)
	}

	return creators, calculateMetadata(totalCount, filter.Page, filter.PageSize), nil
}

func (m *CreatorModel) GetCount(filter *Filter) (int, error) {
	query := `
			SELECT COUNT(*)
			FROM (
				SELECT c.id, c.name, c.slug, c.page_url, c.profile_img, c.date_established
				FROM creator as c
				WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {

		query += fmt.Sprintf(`AND
            to_tsvector('english', c.name || ' ' || COALESCE(c.alias, '')) @@ websearch_to_tsquery('english', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	query += " ) as grouped_count"

	var count int
	err := m.DB.QueryRow(context.Background(), query, args...).Scan(&count)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}

	return count, nil
}

func (m *CreatorModel) GetCreatorRefs(ids []int) ([]*CreatorRef, error) {
	if len(ids) < 1 {
		return nil, nil
	}

	stmt := `SELECT id, name, slug, profile_img
			FROM creator
			WHERE id IN (%s)`

	args := []any{}
	queryPlaceholders := []string{}
	for i, id := range ids {
		queryPlaceholders = append(queryPlaceholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	stmt = fmt.Sprintf(stmt, strings.Join(queryPlaceholders, ","))
	rows, err := m.DB.Query(context.Background(), stmt, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var creators []*CreatorRef
	for rows.Next() {
		c := CreatorRef{}
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.ProfileImage)
		if err != nil {
			return nil, err
		}
		creators = append(creators, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return creators, nil
}

func (m *CreatorModel) Insert(creator *Creator) (int, error) {
	stmt := `
	INSERT INTO creator (name, page_url, date_established, slug, profile_img, alias)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id;
	`

	var id int
	row := m.DB.QueryRow(
		context.Background(), stmt, creator.Name, creator.URL,
		creator.EstablishedDate, creator.Slug, creator.ProfileImage,
		creator.Alias)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	creator.ID = &id

	return id, err
}

func (m *CreatorModel) Search(query string) ([]*Creator, error) {
	query = "%" + query + "%"
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

func (m *CreatorModel) GetById(id int) (*Creator, error) {
	stmt := `SELECT id, name, alias, slug, page_url, profile_img, date_established 
	FROM creator
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.Alias, &c.Slug,
		&c.URL, &c.ProfileImage, &c.EstablishedDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil
}

func (m *CreatorModel) GetCast(id int) ([]*Person, error) {
	stmt := `
		SELECT DISTINCT p.id, p.first, p.last, p.profile_img, p.birthdate, p.slug
		FROM person as p
		JOIN cast_members as cm ON p.id = cm.person_id 
		JOIN sketch as v ON cm.sketch_id = v.id
		JOIN sketch_creator_rel as scr ON v.id = scr.sketch_id
		JOIN creator as c ON scr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		WHERE c.id = $1
		AND cm.role = 'cast'
		AND sh.id is null;
	`

	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}

	var people []*Person
	for rows.Next() {
		var p Person
		destinations := []any{
			&p.ID, &p.First, &p.Last, &p.ProfileImg, &p.BirthDate, &p.Slug,
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}

		people = append(people, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return people, nil

}

func (m *CreatorModel) Exists(id int) (bool, error) {
	stmt := `SELECT EXISTS(
		SELECT 1 FROM creator WHERE id = $1
	)`
	var exists bool
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&exists)
	if err == pgx.ErrNoRows {
		err = nil
	}

	return exists, err
}

func (m *CreatorModel) VectorSearch(query string) ([]*ProfileResult, error) {
	stmt := `
		SELECT id, name, profile_img, slug, date_established, ts_rank(search_vector, websearch_to_tsquery('english', $1)) AS rank
		FROM creator
		WHERE search_vector @@ websearch_to_tsquery('english', $1)
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

func (m *CreatorModel) Update(creator *Creator) error {
	stmt := `
		UPDATE creator SET name = $1, page_url = $2, 
		date_established = $3, profile_img = $4, slug = $5, alias = $6
		WHERE id = $7
	`
	_, err := m.DB.Exec(
		context.Background(), stmt, creator.Name,
		creator.URL, creator.EstablishedDate, creator.ProfileImage,
		creator.Slug, creator.Alias, creator.ID,
	)
	return err
}

func (m *CreatorModel) SearchCount(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM creator
		WHERE search_vector @@ websearch_to_tsquery('english', $1)
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
