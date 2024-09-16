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
	URL             string
	ProfileImage    string
	EstablishedDate time.Time
}

type CreatorModelInterface interface {
	Insert(name, url, imgName, imgExt string, establishedDate time.Time) (int, string, error)
	Get(id int) (*Creator, error)
	ExistsByName(name string) (int, error)
}

type CreatorModel struct {
	DB *pgxpool.Pool
}

func (m *CreatorModel) Insert(name, url, imgName, imgExt string, establishedDate time.Time) (int, string, error) {
	stmt := `
	INSERT INTO creator (name, page_url, date_established, profile_img_path)
	VALUES ($1,$2,$3,CONCAT($4::text, '-', currval(pg_get_serial_sequence('creator', 'id')), $5::text)) 
	RETURNING id, profile_img_path;`

	var id int
	var fullImgName string
	row := m.DB.QueryRow(context.Background(), stmt, name, url, establishedDate, imgName, imgExt)
	err := row.Scan(&id, &fullImgName)
	if err != nil {
		return 0, "", err
	}
	return id, fullImgName, err
}

func (m *CreatorModel) Get(id int) (*Creator, error) {
	stmt := `SELECT id, name, profile_img, creation_date FROM creator
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.ProfileImage, &c.EstablishedDate)
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
