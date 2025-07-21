package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Series struct {
	ID            *int
	Slug          *string
	Title         *string
	Description   *string
	ThumbnailName *string
	Sketches      []*Sketch
}

type SeriesModelInterface interface {
	Delete(id int) error
	GetById(id int) (*Series, error)
	Insert(*Series) (int, error)
	Search(string) ([]*Series, error)
	Update(*Series) error
}

type SeriesModel struct {
	DB *pgxpool.Pool
}

func (m *SeriesModel) Delete(id int) error {
	stmt := `
		DELETE FROM series where id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *SeriesModel) GetById(id int) (*Series, error) {
	stmt := `
		SELECT s.id, s.slug, s.title, s.description, s.thumbnail_name,
		sk.id, sk.slug, sk.title, sk.sketch_url, sk.thumbnail_name, 
		sk.upload_date, sk.episode_start, sk.sketch_number, sk.part_number,
		e.id, e.slug, e.episode_number, e.title, e.air_date, 
		e.thumbnail_name, e.url, e.youtube_id,
		se.id, se.slug, se.season_number,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.profile_img, sh.slug
		FROM series as s
		LEFT JOIN sketch AS sk ON s.id = sk.series_id
		LEFT JOIN episode as e ON sk.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN sketch_creator_rel as skcr ON sk.id = skcr.sketch_id
		LEFT JOIN creator as c ON skcr.creator_id = c.id
		WHERE s.id = $1
		ORDER BY sk.part_number
	`
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	s := &Series{}
	hasRows := false
	for rows.Next() {
		sk := &Sketch{}
		c := &Creator{}
		sh := &Show{}
		ep := &Episode{}
		se := &Season{}
		hasRows = true
		err := rows.Scan(
			&s.ID, &s.Slug, &s.Title, &s.Description, &s.ThumbnailName,
			&sk.ID, &sk.Slug, &sk.Title, &sk.URL, &sk.ThumbnailName,
			&sk.UploadDate, &sk.EpisodeStart, &sk.Number, &sk.SeriesPart,
			&ep.ID, &ep.Slug, &ep.Number, &ep.Title, &ep.AirDate, &ep.Thumbnail,
			&ep.URL, &ep.YoutubeID,
			&se.ID, &se.Slug, &se.Number,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
		)

		if err != nil {
			return nil, err
		}

		sk.Creator = c
		sk.Show = sh

		s.Sketches = append(s.Sketches, sk)
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SeriesModel) Insert(series *Series) (int, error) {
	stmt := `
		INSERT INTO series (slug, title, thumbnail_name)
		VALUES ($1,$2,$3)
		RETURNING id;
		`

	result := m.DB.QueryRow(
		context.Background(),
		stmt, series.Slug, series.Title, series.ThumbnailName,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SeriesModel) Search(query string) ([]*Series, error) {
	query = "%" + query + "%"
	stmt := `SELECT s.id, s.slug, s.title, s.thumbnail_name
			FROM series as s
			WHERE s.title ILIKE $1
			ORDER BY s.title`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	series := []*Series{}
	for rows.Next() {
		s := &Series{}
		err := rows.Scan(
			&s.ID, &s.Slug, &s.Title, &s.ThumbnailName,
		)
		if err != nil {
			return nil, err
		}

		series = append(series, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return series, nil
}

func (m *SeriesModel) Update(series *Series) error {
	stmt := `
		UPDATE series 
		SET slug = $1, title = $2, thumbnail_name = $3
		WHERE id = $4
	`

	_, err := m.DB.Exec(
		context.Background(),
		stmt, series.Slug, series.Title,
		series.ThumbnailName, series.ID,
	)

	return err
}
