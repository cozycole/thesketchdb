package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Recurring struct {
	ID            *int
	Slug          *string
	Title         *string
	Description   *string
	ThumbnailName *string
	Sketches      []*SketchRef
}

type RecurringModelInterface interface {
	Delete(id int) error
	GetById(id int) (*Recurring, error)
	Insert(*Recurring) (int, error)
	Search(string) ([]*Recurring, error)
	Update(*Recurring) error
}

type RecurringModel struct {
	DB *pgxpool.Pool
}

func (m *RecurringModel) Delete(id int) error {
	stmt := `
		DELETE FROM recurring where id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *RecurringModel) GetById(id int) (*Recurring, error) {
	stmt := `
		SELECT r.id, r.slug, r.title, r.description, r.thumbnail_name,
		sk.id, sk.slug, sk.title, sk.thumbnail_name, 
		sk.upload_date, sk.sketch_number, sk.part_number,
		e.id, e.slug, e.episode_number, e.air_date, e.thumbnail_name,
		se.id, se.slug, se.season_number,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.profile_img, sh.slug
		FROM recurring as r
		LEFT JOIN sketch AS sk ON r.id = sk.recurring_id
		LEFT JOIN episode as e ON sk.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN sketch_creator_rel as skcr ON sk.id = skcr.sketch_id
		LEFT JOIN creator as c ON skcr.creator_id = c.id
		WHERE r.id = $1
		ORDER BY sk.upload_date asc
	`
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	s := &Recurring{}
	hasRows := false
	for rows.Next() {
		sk := &SketchRef{}
		c := &CreatorRef{}
		sh := &ShowRef{}
		ep := &EpisodeRef{}
		se := &SeasonRef{}
		hasRows = true
		err := rows.Scan(
			&s.ID, &s.Slug, &s.Title, &s.Description, &s.ThumbnailName,
			&sk.ID, &sk.Slug, &sk.Title, &sk.Thumbnail,
			&sk.UploadDate, &sk.Number,
			&ep.ID, &ep.Slug, &ep.Number, &ep.AirDate,
			&se.ID, &se.Slug, &se.Number,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
		)

		if err != nil {
			return nil, err
		}
		ep.Season = se
		se.Show = sh

		sk.Creator = c
		sk.Episode = ep

		if sk.ID != nil {
			s.Sketches = append(s.Sketches, sk)
		}
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *RecurringModel) Insert(recurring *Recurring) (int, error) {
	stmt := `
		INSERT INTO recurring (slug, title, description, thumbnail_name)
		VALUES ($1,$2,$3,$4)
		RETURNING id;
		`

	result := m.DB.QueryRow(
		context.Background(),
		stmt, recurring.Slug, recurring.Title,
		recurring.Description, recurring.ThumbnailName,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *RecurringModel) Search(query string) ([]*Recurring, error) {
	query = "%" + query + "%"
	stmt := `SELECT r.id, r.slug, r.title, r.thumbnail_name
			FROM recurring as r
			WHERE r.title ILIKE $1
			ORDER BY r.title`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recurring := []*Recurring{}
	for rows.Next() {
		s := &Recurring{}
		err := rows.Scan(
			&s.ID, &s.Slug, &s.Title, &s.ThumbnailName,
		)
		if err != nil {
			return nil, err
		}

		recurring = append(recurring, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return recurring, nil
}

func (m *RecurringModel) Update(recurring *Recurring) error {
	stmt := `
		UPDATE recurring 
		SET slug = $1, title = $2, 
		description =$3, thumbnail_name = $4
		WHERE id = $5
	`

	_, err := m.DB.Exec(
		context.Background(),
		stmt, recurring.Slug, recurring.Title,
		recurring.Description, recurring.ThumbnailName,
		recurring.ID,
	)

	return err
}
