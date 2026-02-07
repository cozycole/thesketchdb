package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Series struct {
	ID            *int
	Slug          *string
	Title         *string
	Description   *string
	ThumbnailName *string
	Sketches      []*SketchRef
}

type SeriesRef struct {
	ID            *int    `json:"id"`
	Slug          *string `json:"slug"`
	Title         *string `json:"title"`
	ThumbnailName *string `json:"thumbnailName"`
}

type SeriesModelInterface interface {
	Delete(id int) error
	GetById(id int) (*Series, error)
	Insert(*Series) (int, error)
	List(f *Filter) ([]*SeriesRef, Metadata, error)
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
		sk.id, sk.slug, sk.title, sk.thumbnail_name, 
		sk.upload_date, sk.sketch_number, sk.part_number,
		e.id, e.slug, e.episode_number, e.air_date, 
		e.thumbnail_name,
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
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
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

func (m *SeriesModel) List(f *Filter) ([]*SeriesRef, Metadata, error) {
	query := `SELECT count(*) OVER(), s.id, s.slug, s.title, s.thumbnail_name%s
			FROM series as s
			WHERE 1=1
	`

	args := []any{}
	argIndex := 1
	if f.Query != "" {
		rankParam := fmt.Sprintf(`
		, ts_rank(
			setweight(to_tsvector('english', s.title) , 'A'),
			websearch_to_tsquery('english', $%d)
		) AS rank
		`, argIndex)

		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`AND
			(
				to_tsvector('english', s.title) @@ websearch_to_tsquery('english', $%d)
				OR
				s.title ILIKE '%%' || $%d || '%%'
			)
		`, argIndex, argIndex)

		args = append(args, f.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := m.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	series := []*SeriesRef{}
	var totalCount int
	for rows.Next() {
		var r SeriesRef
		destinations := []any{
			&totalCount, &r.ID, &r.Slug, &r.Title, &r.ThumbnailName,
		}

		var rank *float32
		if f.Query != "" {
			destinations = append(destinations, &rank)
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, Metadata{}, err
		}

		series = append(series, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return series, calculateMetadata(totalCount, f.Page, f.PageSize), nil
}

func (m *SeriesModel) Insert(series *Series) (int, error) {
	stmt := `
		INSERT INTO series (slug, title, description, thumbnail_name)
		VALUES ($1,$2,$3,$4)
		RETURNING id;
		`

	result := m.DB.QueryRow(
		context.Background(),
		stmt, series.Slug, series.Title,
		series.Description, series.ThumbnailName,
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
		SET slug = $1, title = $2, description = $3, thumbnail_name = $4
		WHERE id = $5
	`

	_, err := m.DB.Exec(
		context.Background(),
		stmt, series.Slug, series.Title, series.Description,
		series.ThumbnailName, series.ID,
	)

	return err
}
