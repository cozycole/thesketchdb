package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID            *int
	Name          *string
	Slug          *string
	ParentID      *int
	Subcategories []*Category
	Tags          []*Tag
}

type CategoryInterface interface {
	Exists(id int) (bool, error)
	Get(id int) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	Insert(category *Category) (int, error)
	Search(query string) (*[]*Category, error)
}

type CategoryModel struct {
	DB *pgxpool.Pool
}

func (m *CategoryModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM categories WHERE id = $1`

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

// Gets one layer of sub categories
func (m *CategoryModel) Get(id int) (*Category, error) {
	var c Category
	stmt := `
        SELECT id, name, slug, parent_id 
        FROM categories 
        WHERE id = $1
    `

	err := m.DB.QueryRow(context.Background(), stmt, id).
		Scan(&c.ID, &c.Name, &c.Slug, &c.ParentID)
	if err != nil {
		return nil, err
	}

	stmt = `
        SELECT id, name, slug 
        FROM categories 
        WHERE parent_id = $1
    `
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subCat Category
		if err := rows.Scan(&subCat.ID, &subCat.Name, &subCat.Slug); err != nil {
			return nil, err
		}
		subCat.ParentID = c.ID
		c.Subcategories = append(c.Subcategories, &subCat)
	}

	stmt = `
        SELECT id, name, slug 
        FROM tags 
        WHERE category_id = $1
    `
	tagRows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var t Tag
		if err := tagRows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		c.Tags = append(c.Tags, &t)
	}

	return &c, nil
}

func (m *CategoryModel) GetBySlug(slug string) (*Category, error) {
	stmt := `SELECT id FROM categories WHERE slug = $1`
	var id int
	err := m.DB.QueryRow(context.Background(), stmt, slug).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return m.Get(id)
}

func (m *CategoryModel) Insert(category *Category) (int, error) {
	stmt := `
	INSERT INTO categories (name, slug, parent_id)
	VALUES ($1,$2,$3)
	RETURNING id;
	`
	var id int
	err := m.DB.QueryRow(
		context.Background(), stmt, category.Name, category.Slug, category.ParentID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *CategoryModel) Search(query string) (*[]*Category, error) {
	query = "%" + query + "%"
	stmt := `SELECT c.id, c.slug, c.name, c.parent_id
			FROM categories as c
			WHERE name ILIKE $1
			ORDER BY name`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		c := &Category{}
		err := rows.Scan(
			&c.ID, &c.Slug, &c.Name, &c.ParentID,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &categories, nil
}
