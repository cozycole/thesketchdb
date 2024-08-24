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
