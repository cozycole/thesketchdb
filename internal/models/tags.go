package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tag struct {
	ID       *int
	Name     *string
	Slug     *string
	Category *Category
}

type TagModelInterface interface {
	Exists(id int) (bool, error)
	Get(id int) (*Tag, error)
	// GetBySlug(slug string) (*Tag, error)
	GetByVideo(vidId int) (*[]*Tag, error)
	Insert(category *Tag) (int, error)
	Search(query string) (*[]*Tag, error)
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
		SELECT t.id, t.name, t.slug,
		c.id, c.name, c.slug
		FROM tags as t
		LEFT JOIN categories as c
		ON t.category_id = c.id
		WHERE t.id = $1
	`
	var c Category
	var t Tag
	err := m.DB.QueryRow(context.Background(), stmt, id).
		Scan(&t.ID, &t.Name, &t.Slug, &c.ID, &c.Name, &c.Slug)
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

func (m *TagModel) GetByVideo(id int) (*[]*Tag, error) {
	stmt := `
		SELECT t.id, t.name, t.slug,
		c.id, c.name, c.slug
		FROM tags as t
		LEFT JOIN video_tags as vt
		ON t.id = vt.tag_id 
		LEFT JOIN categories as c
		ON t.category_id = c.id
		WHERE vt.video_id = $1
	`
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &[]*Tag{}, ErrNoRecord
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
	return &tags, nil
}

func (m *TagModel) Insert(tag *Tag) (int, error) {
	stmt := `
	INSERT INTO tags (name, slug, category_id)
	VALUES ($1,$2,$3)
	RETURNING id;
	`
	var id int
	err := m.DB.QueryRow(
		context.Background(), stmt, tag.Name, tag.Slug, tag.Category.ID,
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
			JOIN categories as c
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
