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
	Alias       *string
	Professions *string
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
	Delete(id int) error
	Get(filter *Filter) ([]*Person, error)
	GetById(id int) (*Person, error)
	GetCount(filter *Filter) (int, error)
	GetCreatorShowCounts(id int) ([]*CreatorShowCounts, error)
	GetPeople(ids []int) ([]*Person, error)
	GetPersonStats(id int) (*PersonStats, error)
	Exists(id int) (bool, error)
	Insert(person *Person) (int, error)
	Search(query string) ([]*Person, error)
	SearchCount(query string) (int, error)
	Update(person *Person) error
}

type PersonModel struct {
	DB *pgxpool.Pool
}

func (m *PersonModel) Delete(id int) error {
	stmt := `
		DELETE FROM person
		WHERE id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *PersonModel) GetPersonStats(id int) (*PersonStats, error) {
	stmt := `
		SELECT
		  (SELECT COUNT(DISTINCT sketch_id)
		   FROM cast_members
		   WHERE person_id = $1) AS sketch_count,
		  (SELECT COUNT(*)
		   FROM cast_members as cm
		   JOIN character as c ON cm.character_id = c.id
		   WHERE c.person_id = $1) AS portrayal_count,
		  (SELECT COUNT(DISTINCT c.creator_id)
		   FROM cast_members as cm
		   JOIN sketch as v ON v.id = cm.sketch_id
		   JOIN  sketch_creator_rel as c ON v.id = c.sketch_id
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

func (m *PersonModel) Insert(person *Person) (int, error) {
	stmt := `
	INSERT INTO person (first, last, aliases, birthdate, professions, slug, profile_img)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id;
	`
	var id int
	row := m.DB.QueryRow(
		context.Background(), stmt, person.First, person.Last, person.Alias,
		person.BirthDate, person.Professions, person.Slug, person.ProfileImg,
	)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
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
	query := `SELECT p.id, p.first, p.last, p.aliases, p.profile_img, p.birthdate, p.slug%s
			FROM person as p
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		rankParam := fmt.Sprintf(`
		, ts_rank(
			setweight(to_tsvector('simple', p.first), 'A') ||
			setweight(to_tsvector('simple', p.last), 'A') ||
			setweight(to_tsvector('simple', COALESCE(p.aliases, '')), 'B'),
			websearch_to_tsquery('simple', $%d)
		) AS rank
		`, argIndex)

		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`AND
            to_tsvector('simple', p.first || ' ' ||
            p.last || ' ' || COALESCE(p.aliases, '')
		) @@ websearch_to_tsquery('simple', $%d)
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
			&p.ID, &p.First, &p.Last, &p.Alias,
			&p.ProfileImg, &p.BirthDate, &p.Slug,
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
	stmt := `SELECT id, first, last, aliases, profile_img, birthdate, slug, professions
			FROM person
			WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	p := &Person{}

	err := row.Scan(
		&p.ID, &p.First, &p.Last, &p.Alias, &p.ProfileImg,
		&p.BirthDate, &p.Slug, &p.Professions,
	)
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
            to_tsvector(
			'simple', COALESCE(p.first,'') || ' ' 
				|| COALESCE(p.last,'') || ' ' || 
				COALESCE(p.aliases, '')) @@ websearch_to_tsquery('simple', $%d)
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

type CreatorShowCounts struct {
	Type      *string
	ID        *int
	Slug      *string
	ImageName *string
	Name      *string
	Count     *int
}

func (m *PersonModel) GetCreatorShowCounts(id int) ([]*CreatorShowCounts, error) {
	stmt := `
		SELECT
		  type,
		  id,
		  slug,
		  image,
		  name,
		  COUNT(DISTINCT sketch_id) AS sketch_count
		FROM (
		  SELECT
			'show' AS type,
			s.id as id,
			s.slug as slug,
			s.profile_img as image,
			s.name AS name,
			cm.sketch_id
		  FROM
			cast_members cm
			JOIN sketch sk ON cm.sketch_id = sk.id
			JOIN episode e ON sk.episode_id = e.id
			JOIN season se ON e.season_id = se.id
			JOIN show s ON se.show_id = s.id
		  WHERE
			cm.person_id = $1

		  UNION ALL
		 
		  SELECT
			'creator' AS type,
			c.id as id,
			c.slug as slug,
			c.profile_img as image,
			c.name AS name,
			cm.sketch_id
		  FROM
			cast_members cm
			JOIN sketch sk ON cm.sketch_id = sk.id
			JOIN sketch_creator_rel scr ON sk.id = scr.sketch_id
			JOIN creator c ON scr.creator_id = c.id
		  WHERE
			cm.person_id = $1
		) combined
		GROUP BY
		  type, id, slug, image, name
		ORDER BY
		  sketch_count DESC;
		`

	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := []*CreatorShowCounts{}
	for rows.Next() {
		c := &CreatorShowCounts{}
		err := rows.Scan(
			&c.Type, &c.ID, &c.Slug, &c.ImageName, &c.Name, &c.Count,
		)
		if err != nil {
			return nil, err
		}
		counts = append(counts, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return counts, nil
}

func (m *PersonModel) Search(query string) ([]*Person, error) {
	query = "%" + query + "%"
	stmt := `SELECT id, slug, first, last, profile_img, birthdate
			FROM person
			WHERE CONCAT(LOWER(first), LOWER(last)) LIKE LOWER($1)
			OR LOWER(last) LIKE LOWER($1)
			LIMIT 10`

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

func (m *PersonModel) Update(person *Person) error {
	stmt := `
	UPDATE person SET first = $1, last = $2, professions = $3, 
	profile_img = $4, birthdate = $5, slug = $6, aliases = $7
	WHERE id = $8`
	_, err := m.DB.Exec(
		context.Background(), stmt, person.First,
		person.Last, person.Professions, person.ProfileImg,
		person.BirthDate, person.Slug, person.Alias, person.ID,
	)
	return err
}
