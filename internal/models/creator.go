package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Creator struct {
	ID              int
	Name            string
	ProfileImage    string
	DateEstablished time.Time
}

type CreatorModelInterface interface {
	Get(id int) (*Creator, error)
	ExistsByName(name string) (int, error)
}

type CreatorModel struct {
	DB *pgxpool.Pool
}

func (m *CreatorModel) Get(id int) (*Creator, error) {
	stmt := `SELECT id, name, profile_img, creation_date FROM creator
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.ProfileImage, &c.DateEstablished)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil
}

func (m *CreatorModel) ExistsByName(name string) (int, error) {
	stmt := `SELECT id FROM creator WHERE name = $1`
	row := m.DB.QueryRow(context.Background(), stmt, name)

	c := &Creator{}

	err := row.Scan(&c.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	return c.ID, nil
}
