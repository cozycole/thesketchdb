package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Video struct {
	ID         int
	Title      string
	URL        string
	Slug       string
	Thumbnail  string
	Rating     string
	UploadDate time.Time
	Creator    *Creator
}

type VideoModelInterface interface {
	Search(search string, offset int) ([]*Video, error)
	GetAll(offset int) ([]*Video, error)
	Get(id int) (*Video, error)
	GetBySlug(slug string) (*Video, error)
	GetByCreator(id int) ([]*Video, error)
	Insert(title, video_url, rating, imgName, imgExt string, upload_date time.Time) (int, string, string, error)
	InsertVideoCreatorRelation(vidId, creatorId int) error
	InsertVideoPersonRelation(vidId, personId int) error
}

type VideoModel struct {
	DB         *pgxpool.Pool
	ResultSize int
}

func (m *VideoModel) Search(search string, offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.profile_img
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE search_vector @@ to_tsquery('english', $1)
		LIMIT $2
		OFFSET $3;
	`
	rows, err := m.DB.Query(context.Background(), stmt, search, m.ResultSize, m.ResultSize*offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		err := rows.Scan(&v.ID, &v.Title, &v.URL, &v.Thumbnail)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

func (m *VideoModel) GetBySlug(slug string) (*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.slug = $1
	`
	row := m.DB.QueryRow(context.Background(), stmt, slug)
	v := &Video{}
	c := &Creator{}
	err := row.Scan(
		&v.ID, &v.Title, &v.URL, &v.Slug, &v.Thumbnail, &v.UploadDate,
		&c.ID, &c.Name, &c.URL, &c.Slug, &c.ProfileImage,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	v.Creator = c

	return v, nil
}
func (m *VideoModel) GetByCreator(id int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE c.id = $1
	`
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		c := &Creator{}
		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.Thumbnail, &v.UploadDate,
			&c.ID, &c.Name, &c.URL, &c.Slug, &c.ProfileImage,
		)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil

}

func (m *VideoModel) Get(id int) (*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.id = $1
	`
	row, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer row.Close()
	v := &Video{}
	c := &Creator{}
	err = row.Scan(
		&v.ID, &v.Title, &v.URL, &v.Thumbnail, &v.UploadDate,
		&c.ID, &c.Name, &c.ProfileImage,
	)
	if err != nil {
		return nil, err
	}
	v.Creator = c

	return v, nil
}

// Will make DRY later
func (m *VideoModel) GetAll(offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		LIMIT $1
		OFFSET $2;
	`
	rows, err := m.DB.Query(context.Background(), stmt, m.ResultSize, m.ResultSize*offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		c := &Creator{}
		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.Thumbnail, &v.UploadDate,
			&c.ID, &c.Name, &c.URL, &c.Slug, &c.ProfileImage,
		)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

func (m *VideoModel) Insert(title, video_url, rating, slug, imgExt string, upload_date time.Time) (int, string, string, error) {
	stmt := `
	INSERT INTO video (title, video_url, upload_date, pg_rating, slug, thumbnail_name)
	VALUES ($1,$2,$3,$4,
		CONCAT($5::text, '-', currval(pg_get_serial_sequence('video', 'id'))), 
		CONCAT($5::text, '-', currval(pg_get_serial_sequence('video', 'id')), $6::text))
	RETURNING id, slug, thumbnail_name;`

	result := m.DB.QueryRow(
		context.Background(), stmt, title,
		video_url, upload_date, rating,
		slug, imgExt,
	)

	var id int
	var imgName string
	err := result.Scan(&id, &slug, &imgName)
	if err != nil {
		return 0, "", "", err
	}
	return id, slug, imgName, nil
}

func (m *VideoModel) InsertVideoCreatorRelation(vidId, creatorId int) error {
	stmt := `INSERT INTO video_creator_rel (video_id, creator_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, creatorId)
	return err
}

func (m *VideoModel) InsertVideoPersonRelation(vidId, personId int) error {
	stmt := `INSERT INTO video_person_rel (video_id, person_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, personId)
	return err
}
