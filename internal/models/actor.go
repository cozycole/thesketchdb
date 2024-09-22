package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Actor struct {
	ID         int
	First      string
	Last       string
	ProfileImg string
	BirthDate  time.Time
}

type ActorModelInterface interface {
	Get(id int) (*Actor, error)
	Exists(id int) (bool, error)
	Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, error)
}

type ActorModel struct {
	DB *pgxpool.Pool
}

func (m *ActorModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, error) {
	stmt := `
	INSERT INTO actor (first, last, birthdate, profile_img)
	VALUES ($1,$2,$3,CONCAT($4::text, '-', currval(pg_get_serial_sequence('actor', 'id')), $5::text)) 
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

func (m *ActorModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM actor WHERE id = $1`
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

func (m *ActorModel) Get(id int) (*Actor, error) {
	stmt := `SELECT id, first, profile_img, birthdate FROM actor
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	a := &Actor{}

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
