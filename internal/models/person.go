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
	First       string
	Last        string
	ProfileImg  string
	BirthDate   *time.Time
	Description *string
}

type PersonModelInterface interface {
	Get(id int) (*Person, error)
	Exists(id int) (bool, error)
	Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, error)
}

type PersonModel struct {
	DB *pgxpool.Pool
}

func (m *PersonModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, error) {
	stmt := `
	INSERT INTO person (first, last, birthdate, profile_img)
	VALUES ($1,$2,$3,CONCAT($4::text, '-', currval(pg_get_serial_sequence('person', 'id')), $5::text)) 
	RETURNING id, profile_img;`
	var id int
	var fullImgName string
	row := m.DB.QueryRow(context.Background(), stmt, first, last, birthDate, imgName, imgExt)
	err := row.Scan(&id, &fullImgName)
	if err != nil {
		return 0, "", err
	}
	return id, fullImgName, err
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

func (m *PersonModel) Get(id int) (*Person, error) {
	stmt := `SELECT id, first, profile_img, birthdate FROM person
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	a := &Person{}

	err := row.Scan(&a.ID, &a.First, &a.Last, &a.ProfileImg, &a.BirthDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return a, nil
}
