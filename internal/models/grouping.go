package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Grouping struct {
	ID          *int         `json:"id"`
	Slug        *string      `json:"slug"`
	Title       *string      `json:"title"`
	Description *string      `json:"description"`
	Position    *int         `json:"position"`
	Sketches    []*SketchRef `json:"sketches"`
}

type GroupingModelInterface interface {
	Delete(id int) error
	GetById(id int) (*Grouping, error)
	Insert(*Grouping) (int, error)
	Search(string) ([]*Grouping, error)
	Update(*Grouping) error
}

type GroupingModel struct {
	DB *pgxpool.Pool
}

func (m *GroupingModel) Delete(id int) error {
	stmt := `
		DELETE FROM grouping where id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *GroupingModel) GetById(id int) (*Grouping, error) {
	stmt := `
		SELECT r.id, r.slug, r.title, r.description, 
		sk.id, sk.slug, sk.title, sk.thumbnail_name, 
		sk.upload_date, sk.sketch_number,
		e.id, e.slug, e.episode_number, e.air_date,
		se.id, se.slug, se.season_number,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img
		FROM grouping as r
		LEFT JOIN sketch AS sk ON r.id = sk.grouping_id
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

	s := &Grouping{}
	hasRows := false
	for rows.Next() {
		sk := &SketchRef{}
		c := &CreatorRef{}
		sh := &ShowRef{}
		ep := &EpisodeRef{}
		se := &SeasonRef{}
		hasRows = true
		err := rows.Scan(
			&s.ID, &s.Slug, &s.Title, &s.Description,
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

func (m *GroupingModel) Insert(grouping *Grouping) (int, error) {
	stmt := `
		INSERT INTO grouping (slug, title, description, show_id, creator_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id;
		`

	result := m.DB.QueryRow(
		context.Background(),
		stmt, grouping.Slug, grouping.Title,
		grouping.Description,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *GroupingModel) Search(query string) ([]*Grouping, error) {
	query = "%" + query + "%"
	stmt := `SELECT g.id, g.slug, g.title
			FROM grouping as g
			WHERE g.title ILIKE $1
			ORDER BY g.title`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouping := []*Grouping{}
	for rows.Next() {
		g := &Grouping{}
		err := rows.Scan(
			&g.ID, &g.Slug, &g.Title, &g.Description,
		)
		if err != nil {
			return nil, err
		}

		grouping = append(grouping, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return grouping, nil
}

func (m *GroupingModel) Update(grouping *Grouping) error {
	stmt := `
		UPDATE sketch_grouping 
		SET slug = $1, title = $2, 
		description =$3, 
		WHERE id = $4
	`

	_, err := m.DB.Exec(
		context.Background(),
		stmt, grouping.Slug, grouping.Title,
		grouping.Description,
		grouping.ID,
	)

	return err
}
