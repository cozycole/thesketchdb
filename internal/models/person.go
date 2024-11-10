package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct {
	ID          int
	Slug        string
	First       string
	Last        string
	ProfileImg  string
	BirthDate   *time.Time
	Description *string
}

type PersonModelInterface interface {
	GetBySlug(slug string) (*Person, error)
	Get(id int) (*Person, error)
	Exists(id int) (bool, error)
	Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, string, error)
}

type PersonModel struct {
	DB *pgxpool.Pool
}

func (m *PersonModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, string, error) {
	stmt := `
	INSERT INTO person (first, last, birthdate, slug, profile_img)
	VALUES ($1,$2,$3,
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('person', 'id')))
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('person', 'id')), $5::text))
	RETURNING id, slug, profile_img;`
	var id int
	var fullImgName, slug string
	row := m.DB.QueryRow(context.Background(), stmt, first, last, birthDate, imgName, imgExt)
	err := row.Scan(&id, &fullImgName, &slug)
	if err != nil {
		return 0, "", "", err
	}
	return id, fullImgName, slug, err
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

func (m *PersonModel) GetBySlug(slug string) (*Person, error) {
	person_id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.Get(person_id)

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

func (m *PersonModel) Get(id int) (*Person, error) {
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
