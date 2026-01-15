package models

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tag struct {
	ID       *int
	Name     *string
	Slug     *string
	Type     *string
	Category *Category
	Count    *int
}

type TagRef struct {
	ID       *int
	Name     *string
	Slug     *string
	Category *CategoryRef
}

type TagModelInterface interface {
	Exists(id int) (bool, error)
	Get(id int) (*Tag, error)
	GetTags(ids []int) ([]*Tag, error)
	GetTagRefs(ids []int) ([]*TagRef, error)
	GetTagsByType(string) ([]*Tag, error)
	GetBySketch(sketchId int) ([]*Tag, error)
	Insert(category *Tag) (int, error)
	Search(query string) (*[]*Tag, error)
	Update(*Tag) error
}

type TagModel struct {
	DB *pgxpool.Pool
}

func (m *TagModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM tags WHERE id = $1`

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

func (m *TagModel) Get(id int) (*Tag, error) {
	stmt := `
		SELECT t.id, t.name, t.slug, t.type,
		c.id, c.name, c.slug
		FROM tags as t
		LEFT JOIN categories as c
		ON t.category_id = c.id
		WHERE t.id = $1
	`
	var c Category
	var t Tag
	err := m.DB.QueryRow(context.Background(), stmt, id).
		Scan(&t.ID, &t.Name, &t.Slug, &t.Type, &c.ID, &c.Name, &c.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	t.Category = &c
	return &t, nil
}

func (m *TagModel) GetTags(ids []int) ([]*Tag, error) {
	if ids != nil && len(ids) < 1 {
		return nil, nil
	}

	stmt := `SELECT t.id, t.name, c.id, c.name
			FROM tags as t
			LEFT JOIN categories as c ON t.category_id = c.id
			WHERE t.id IN (%s)`

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

	var tags []*Tag
	for rows.Next() {
		t := Tag{}
		c := Category{}

		err := rows.Scan(&t.ID, &t.Name, &c.ID, &c.Name)
		if err != nil {
			return nil, err
		}

		t.Category = &c
		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) GetTagRefs(ids []int) ([]*TagRef, error) {
	if len(ids) < 1 {
		return nil, nil
	}

	stmt := `SELECT t.id, t.slug, t.name, c.id, c.slug, c.name
			FROM tags as t
			LEFT JOIN categories as c ON t.category_id = c.id
			WHERE t.id IN (%s)`

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

	var tags []*TagRef
	for rows.Next() {
		t := TagRef{}
		c := CategoryRef{}

		err := rows.Scan(&t.ID, &t.Slug, &t.Name, &c.ID, &c.Slug, &c.Name)
		if err != nil {
			return nil, err
		}

		t.Category = &c
		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) GetBySketch(id int) ([]*Tag, error) {
	stmt := `
		SELECT t.id, t.name, t.slug,
		c.id, c.name, c.slug
		FROM tags as t
		LEFT JOIN sketch_tags as vt
		ON t.id = vt.tag_id 
		LEFT JOIN categories as c
		ON t.category_id = c.id
		WHERE vt.sketch_id = $1
		ORDER by t.id asc
	`
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*Tag{}, ErrNoRecord
		} else {
			return nil, err
		}
	}

	var tags []*Tag
	for rows.Next() {
		var c Category
		var t Tag
		err := rows.Scan(&t.ID, &t.Name, &t.Slug, &c.ID, &c.Name, &c.Slug)
		if err != nil {
			return nil, err
		}

		t.Category = &c

		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (m *TagModel) GetTagsByType(tagType string) ([]*Tag, error) {
	stmt := `
		SELECT t.id, t.name, t.slug,
		c.id, c.name, c.slug
		FROM tags as t
		LEFT JOIN categories as c ON t.category_id = c.id
		WHERE t.type = $1
	`

	rows, err := m.DB.Query(context.Background(), stmt, tagType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*Tag{}, ErrNoRecord
		} else {
			return nil, err
		}
	}

	var tags []*Tag
	for rows.Next() {
		var c Category
		var t Tag
		err := rows.Scan(&t.ID, &t.Name, &t.Slug, &c.ID, &c.Name, &c.Slug)
		if err != nil {
			return nil, err
		}

		t.Category = &c

		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (m *TagModel) Insert(tag *Tag) (int, error) {
	stmt := `
	INSERT INTO tags (name, slug, type, category_id)
	VALUES ($1,$2,$3,$4)
	RETURNING id;
	`
	var categoryId *int
	if tag.Category != nil &&
		safeDeref(tag.Category.ID) != 0 {
		categoryId = tag.Category.ID
	}

	var id int
	err := m.DB.QueryRow(
		context.Background(), stmt, tag.Name, tag.Slug,
		tag.Type, categoryId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *TagModel) Search(query string) (*[]*Tag, error) {
	query = "%" + query + "%"
	stmt := `SELECT t.id, t.slug, t.name, 
			c.id, c.slug, c.name
			FROM tags as t
			LEFT JOIN categories as c
			ON t.category_id = c.id
			WHERE t.name ILIKE $1
			ORDER BY c.name, t.name`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []*Tag{}
	for rows.Next() {
		c := &Category{}
		t := &Tag{}
		err := rows.Scan(
			&t.ID, &t.Slug, &t.Name,
			&c.ID, &c.Slug, &c.Name,
		)
		if err != nil {
			return nil, err
		}

		t.Category = c
		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &tags, nil
}

func (m *TagModel) Update(tag *Tag) error {
	stmt := `
		UPDATE tags SET name = $1, slug = $2, type = $3, category_id = $4
		WHERE id = $5
	`
	var categoryId *int
	if tag.Category != nil && safeDeref(tag.Category.ID) != 0 {
		categoryId = new(int)
		*categoryId = safeDeref(tag.Category.ID)
	}
	_, err := m.DB.Exec(
		context.Background(), stmt, tag.Name, tag.Slug,
		tag.Type, categoryId, tag.ID,
	)
	return err
}
