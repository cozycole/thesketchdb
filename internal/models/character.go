package models

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Character struct {
	ID          *int
	Slug        *string
	Name        *string
	Aliases     *string
	Type        *string
	Image       *string
	Description *string
	Portrayal   *Person
}

type CharacterModelInterface interface {
	Delete(id int) error
	Exists(id int) (bool, error)
	Get(filter *Filter) ([]*Character, error)
	GetById(id int) (*Character, error)
	GetCharacters(ids []int) ([]*Character, error)
	GetCount(filter *Filter) (int, error)
	Insert(character *Character) (int, error)
	Search(search string) ([]*Character, error)
	SearchCount(query string) (int, error)
	Update(character *Character) error
	VectorSearch(query string, limit, offset int) ([]*ProfileResult, error)
}

type CharacterModel struct {
	DB *pgxpool.Pool
}

func (m *CharacterModel) Get(filter *Filter) ([]*Character, error) {
	query := `SELECT c.id, c.slug, c.name, c.img_name,
			p.id, p.slug, p.first, p.last, p.profile_img %s
			FROM character as c
			LEFT JOIN person as p ON c.person_id = p.id
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		rankParam := fmt.Sprintf(`
		, ts_rank(
			setweight(to_tsvector('english', c.name), 'A') ||
			setweight(to_tsvector('english', c.aliases), 'A') ||
			setweight(to_tsvector('simple', p.first), 'B') ||
			setweight(to_tsvector('simple', p.last), 'B'),
			websearch_to_tsquery('english', $%d)
		) AS rank
		`, argIndex)

		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`AND
            to_tsvector('english', COALESCE(c.name, '') || ' ' || COALESCE(c.aliases, '') || 
			' ' || COALESCE(p.first,'') || ' ' || COALESCE(p.last,'')) @@ websearch_to_tsquery('english', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	fmt.Println(query)

	rows, err := m.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	var characters []*Character
	for rows.Next() {
		var c Character
		var p Person
		destinations := []any{
			&c.ID, &c.Slug, &c.Name, &c.Image,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
		}

		var rank *float32
		if filter.Query != "" {
			destinations = append(destinations, &rank)
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}

		characters = append(characters, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

func (m *CharacterModel) GetCharacters(ids []int) ([]*Character, error) {
	if len(ids) < 1 {
		return nil, nil
	}

	stmt := `SELECT id, slug, name, img_name 
			FROM character
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

	var characters []*Character
	for rows.Next() {
		c := Character{}
		err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Image)
		if err != nil {
			return nil, err
		}
		characters = append(characters, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

func (m *CharacterModel) GetById(id int) (*Character, error) {
	stmt := `SELECT c.id, c.slug, c.name, c.aliases, c.character_type, c.img_name,
			p.id, p.slug, p.first, p.last
			FROM character AS c
			LEFT JOIN person AS p ON c.person_id = p.id
			WHERE c.id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	c := &Character{}
	p := &Person{}

	err := row.Scan(&c.ID, &c.Slug, &c.Name, &c.Aliases, &c.Type, &c.Image,
		&p.ID, &p.Slug, &p.First, &p.Last)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	if p.ID != nil {
		c.Portrayal = p
	}

	return c, nil
}

func (m *CharacterModel) GetCount(filter *Filter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM (
			SELECT c.id, c.slug, c.name, c.img_name,
			p.id, p.slug, p.first, p.last, p.profile_img, p.birthdate
			FROM character as c
			LEFT JOIN person as p ON c.person_id = p.id
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		query += fmt.Sprintf(`AND
            to_tsvector('english', c.name || ' ' || COALESCE(c.aliases, '') || ' ' ||
			COALESCE(p.first,'') || ' ' || COALESCE(p.last,'')) @@ websearch_to_tsquery('english', $%d)
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

func (m *CharacterModel) Delete(id int) error {
	stmt := `
		DELETE FROM character
		WHERE id = $1
	`
	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *CharacterModel) Insert(character *Character) (int, error) {
	stmt := `
	INSERT INTO character (name, aliases, character_type, slug, img_name, person_id)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id;
	`
	var personId *int
	if character.Portrayal != nil &&
		character.Portrayal.ID != nil &&
		*character.Portrayal.ID != 0 {
		personId = character.Portrayal.ID
	}

	var id int
	row := m.DB.QueryRow(
		context.Background(), stmt, character.Name,
		character.Aliases, character.Type, character.Slug,
		character.Image, personId,
	)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
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

func (m *CharacterModel) Update(character *Character) error {
	stmt := `
		UPDATE character SET name = $1, aliases = $2, character_type = $3, img_name = $4, 
		person_id = $5, slug = $6
		WHERE id = $7
	`
	var personId *int
	if character.Portrayal != nil &&
		character.Portrayal.ID != nil &&
		*character.Portrayal.ID != 0 {
		personId = character.Portrayal.ID
	}

	_, err := m.DB.Exec(
		context.Background(), stmt, character.Name,
		character.Aliases, character.Type, character.Image,
		personId, character.Slug, character.ID,
	)
	return err
}

func (m *CharacterModel) VectorSearch(query string, limit, offset int) ([]*ProfileResult, error) {
	fmt.Printf("Got here %s %d %d\n", query, limit, offset)
	stmt := `
		SELECT id, name, img_name, slug, ts_rank(search_vector, websearch_to_tsquery('english', $1)) AS rank
		FROM character
		WHERE search_vector @@ websearch_to_tsquery('english', $1)
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
		WHERE c.search_vector @@ websearch_to_tsquery('english', $1)
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
