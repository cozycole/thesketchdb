package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct {
	ID          *int
	Slug        *string
	First       *string
	Last        *string
	ProfileImg  *string
	BirthDate   *time.Time
	Description *string
}

type PersonStats struct {
	SketchCount    int
	CharacterCount int
	CreatorCount   int
	PortrayalCount int
}

type PersonModelInterface interface {
	GetBySlug(slug string) (*Person, error)
	Get(filter *Filter) ([]*Person, error)
	GetById(id int) (*Person, error)
	GetCount(filter *Filter) (int, error)
	GetPeople(ids []int) ([]*Person, error)
	GetPersonStats(id int) (*PersonStats, error)
	Exists(id int) (bool, error)
	Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, string, error)
	Search(query string) ([]*Person, error)
	VectorSearch(query string) ([]*ProfileResult, error)
	SearchCount(query string) (int, error)
}

type PersonModel struct {
	DB *pgxpool.Pool
}

func (m *PersonModel) GetPersonStats(id int) (*PersonStats, error) {
	stmt := `
		SELECT
		  (SELECT COUNT(DISTINCT video_id)
		   FROM cast_members
		   WHERE person_id = $1) AS sketch_count,
		  (SELECT COUNT(*)
		   FROM cast_members as cm
		   JOIN character as c ON cm.character_id = c.id
		   WHERE c.person_id = $1) AS portrayal_count,
		  (SELECT COUNT(DISTINCT c.creator_id)
		   FROM cast_members as cm
		   JOIN video as v ON v.id = cm.video_id
		   JOIN  video_creator_rel as c ON v.id = c.video_id
		   WHERE cm.person_id = $1) AS creator_count,
		  (SELECT COUNT(DISTINCT cm.character_id)
		   FROM cast_members cm
		   WHERE cm.person_id = $1 AND cm.character_id IS NOT NULL) AS character_count;
	`

	stats := &PersonStats{}
	row := m.DB.QueryRow(context.Background(), stmt, id)
	err := row.Scan(&stats.SketchCount, &stats.PortrayalCount, &stats.CreatorCount, &stats.CharacterCount)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (m *PersonModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, string, error) {
	stmt := `
	INSERT INTO person (first, last, birthdate, slug, profile_img)
	VALUES ($1,$2,$3,
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('person', 'id'))),
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('person', 'id')), $5::text))
	RETURNING id, slug, profile_img;`
	var id int
	var fullImgName, slug string
	row := m.DB.QueryRow(context.Background(), stmt, first, last, birthDate, imgName, imgExt)
	err := row.Scan(&id, &slug, &fullImgName)
	if err != nil {
		return 0, "", "", err
	}
	return id, slug, fullImgName, err
}

func (m *PersonModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM person WHERE id = $1`
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

func (m *PersonModel) Get(filter *Filter) ([]*Person, error) {
	query := `SELECT p.id, p.first, p.last, p.profile_img, p.birthdate, p.slug%s
			FROM person as p
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		rankParam := fmt.Sprintf(`
		, ts_rank(
			setweight(to_tsvector('simple', p.first), 'A') ||
			setweight(to_tsvector('simple', p.last), 'A'),
			to_tsquery('simple', $%d)
		) AS rank
		`, argIndex)

		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`AND
            to_tsvector('simple', p.first || ' ' ||
            p.last
		) @@ to_tsquery('simple', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := m.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	var people []*Person
	for rows.Next() {
		var p Person
		destinations := []any{
			&p.ID, &p.First, &p.Last, &p.ProfileImg, &p.BirthDate, &p.Slug,
		}

		var rank *float32
		if filter.Query != "" {
			destinations = append(destinations, &rank)
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

func (m *PersonModel) GetBySlug(slug string) (*Person, error) {
	person_id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.GetById(person_id)

}

func (m *PersonModel) GetIdBySlug(slug string) (int, error) {
	stmt := `SELECT p.id FROM person AS p WHERE p.slug = $1`
	id_row := m.DB.QueryRow(context.Background(), stmt, slug)

	var id int
	err := id_row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *PersonModel) GetById(id int) (*Person, error) {
	stmt := `SELECT id, first, last, profile_img, birthdate, slug 
			FROM person
			WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	p := &Person{}

	err := row.Scan(&p.ID, &p.First, &p.Last, &p.ProfileImg, &p.BirthDate, &p.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *PersonModel) GetPeople(ids []int) ([]*Person, error) {
	if len(ids) < 1 {
		return nil, nil
	}

	stmt := `SELECT id, first, last, profile_img, birthdate, slug 
			FROM person
			WHERE id IN (%s)`

	args := []interface{}{}
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

	var people []*Person
	for rows.Next() {
		p := Person{}
		err := rows.Scan(&p.ID, &p.First, &p.Last, &p.ProfileImg, &p.BirthDate, &p.Slug)
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

func (m *PersonModel) GetCount(filter *Filter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM (
			SELECT p.id, p.first, p.last, p.profile_img, p.birthdate, p.slug
			FROM person as p
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		query += fmt.Sprintf(`AND
            to_tsvector('simple', COALESCE(p.first,'') || ' ' || COALESCE(p.last,'')) @@ to_tsquery('simple', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
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

func (m *PersonModel) Search(query string) ([]*Person, error) {
	query = query + "%"
	stmt := `SELECT id, slug, first, last, profile_img, birthdate
			FROM person
			WHERE CONCAT(LOWER(first), LOWER(last)) LIKE LOWER($1)
			OR LOWER(last) LIKE LOWER($1)`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	people := []*Person{}
	for rows.Next() {
		p := &Person{}
		err := rows.Scan(
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg, &p.BirthDate,
		)
		if err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return people, nil
}

func (m *PersonModel) SearchCount(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM person
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

func (m *PersonModel) VectorSearch(query string) ([]*ProfileResult, error) {
	stmt := `
		SELECT id, first, last, profile_img, birthdate, slug, ts_rank(search_vector, websearch_to_tsquery('english', $1)) AS rank
		FROM person
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
		var first, last *string
		err := rows.Scan(
			&pr.ID, &first, &last, &pr.Img, &pr.Date, &pr.Slug, &pr.Rank,
		)
		if err != nil {
			return nil, err
		}

		name := *first + " " + *last
		resType := "person"
		pr.Name = &name
		pr.Type = &resType

		results = append(results, pr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
