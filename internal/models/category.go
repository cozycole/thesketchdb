package models

import (
	"context"
	"errors"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID   *int
	Name *string
	Slug *string
	Tags []*Tag
}

type CategoryInterface interface {
	Exists(id int) (bool, error)
	Get(id int) (*Category, error)
	GetAll() ([]*Category, error)
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

func (m *CategoryModel) Get(id int) (*Category, error) {
	stmt := `
        SELECT DISTINCT c.id, c.name, t.id, t.name
        FROM categories as c
				LEFT JOIN tags as t ON c.id = t.category_id
				WHERE c.id = $1
    `
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var c Category
	for rows.Next() {
		var t Tag
		if c.ID == nil {
			if err := rows.Scan(&c.ID, &c.Name, &t.ID, &t.Name); err != nil {
				return nil, err
			}
		} else {
			if err := rows.Scan(nil, nil, &t.ID, &t.Name); err != nil {
				return nil, err
			}
		}

		c.Tags = append(c.Tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *CategoryModel) GetAll() ([]*Category, error) {
	stmt := `
			SELECT DISTINCT c.id, c.name, t.id, t.name
			FROM categories as c
			LEFT JOIN tags as t ON c.id = t.category_id
			ORDER BY c.name DESC
    `
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryMap := make(map[int]*Category)
	for rows.Next() {
		var c Category
		var t Tag
		if err := rows.Scan(&c.ID, &c.Name, &t.ID, &t.Name); err != nil {
			return nil, err
		}

		if _, ok := categoryMap[*c.ID]; !ok {
			categoryMap[*c.ID] = &c
		}

		if t.ID != nil {
			mappedCategory := categoryMap[*c.ID]
			mappedCategory.Tags = append(mappedCategory.Tags, &t)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var categories []*Category
	for _, c := range categoryMap {
		categories = append(categories, c)
	}

	// Sort it alphabetically
	sort.Slice(categories, func(i, j int) bool {
		return *categories[i].Name < *categories[j].Name
	})

	return categories, nil
}

func (m *CategoryModel) Insert(category *Category) (int, error) {
	stmt := `
	INSERT INTO categories (name, slug)
	VALUES ($1,$2,$3)
	RETURNING id;
	`
	var id int
	err := m.DB.QueryRow(
		context.Background(), stmt, category.Name, category.Slug,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *CategoryModel) Search(query string) (*[]*Category, error) {
	query = "%" + query + "%"
	stmt := `SELECT c.id, c.slug, c.name,
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
			&c.ID, &c.Slug, &c.Name,
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
