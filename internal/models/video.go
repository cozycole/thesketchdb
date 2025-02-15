package models

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CastMember struct {
	Position      *int
	Actor         *Person
	Character     *Character
	ThumbnailName *string
	ThumbnailFile *multipart.FileHeader
	// TODO: add attributes about the role
}

type Video struct {
	ID            int
	Title         string
	URL           *string
	YoutubeID     *string
	Slug          string
	ThumbnailName string
	ThumbnailFile *multipart.FileHeader
	Rating        string
	Description   *string
	UploadDate    *time.Time
	Creator       *Creator
	Cast          []*CastMember
	Liked         bool
}

type VideoModelInterface interface {
	Get(id int) (*Video, error)
	GetAll(limit int) ([]*Video, error)
	GetByCreator(id int) ([]*Video, error)
	GetByPerson(id int) ([]*Video, error)
	GetBySlug(slug string) (*Video, error)
	GetByUserLikes(id int) ([]*Video, error)
	GetLatest(limit, offset int) ([]*Video, error)
	HasLike(vidId, userId int) (bool, error)
	Insert(video *Video) (int, error)
	InsertThumbnailName(vidId int, name string) error
	InsertVideoCreatorRelation(vidId, creatorId int) error
	InsertVideoPersonRelation(vidId, personId, position int, characterId *int, imgName string) error
	Search(search string, limit, offset int) ([]*Video, error)
	SearchCount(query string) (int, error)
	Update(video *Video) error
}

type VideoModel struct {
	DB *pgxpool.Pool
}

func (m *VideoModel) Get(id int) (*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.youtube_id, v.thumbnail_name, v.upload_date, v.description, v.pg_rating,
			c.id, c.name, c.slug, c.profile_img
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.id = $1
	`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	v := &Video{}
	c := &Creator{}
	err := row.Scan(
		&v.ID, &v.Title, &v.URL, &v.YoutubeID, &v.ThumbnailName, &v.UploadDate, &v.Description, &v.Rating,
		&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	v.Creator = c

	members, err := m.GetCastMembers(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	v.Cast = members

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
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
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

func (m *VideoModel) GetBySlug(slug string) (*Video, error) {
	id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.Get(id)
}

func (m *VideoModel) GetByPerson(id int) ([]*Video, error) {
	stmt := `SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		JOIN video_person_rel as vpl
		ON v.id = vpl.video_id
		JOIN person as p
		ON vpl.person_id = p.id
		WHERE p.id = $1`

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
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
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

func (m *VideoModel) GetIdBySlug(slug string) (int, error) {
	stmt := `SELECT v.id FROM video as v WHERE v.slug = $1`
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

func (m *VideoModel) GetByUserLikes(userId int) ([]*Video, error) {
	stmt := `SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		JOIN likes as l
		ON v.id = l.video_id
		JOIN users as u
		ON l.user_id = u.id
		WHERE u.id = $1
		ORDER BY l.created_at desc
	`

	rows, err := m.DB.Query(context.Background(), stmt, userId)
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
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
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

func (m *VideoModel) GetLatest(limit, offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.upload_date is NOT NULL
		ORDER BY v.upload_date DESC
		LIMIT $1
		OFFSET $2;
	`

	rows, err := m.DB.Query(context.Background(), stmt, limit, offset)
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
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
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

func (m *VideoModel) GetCastMembers(video_id int) ([]*CastMember, error) {
	stmt := `
		SELECT p.id, p.slug, p.first, p.last, p.birthdate,
			p.description, p.profile_img, 
			vpr.position, vpr.img,
			ch.id, ch.name, ch.img_name
		FROM video AS v
		LEFT JOIN video_person_rel as vpr
		ON v.id = vpr.video_id
		LEFT JOIN person as p
		ON vpr.person_id = p.id
		LEFT JOIN character as ch
		ON vpr.character_id = ch.id
		WHERE v.id = $1
	`
	rows, err := m.DB.Query(context.Background(), stmt, video_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	members := []*CastMember{}
	for rows.Next() {
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		err := rows.Scan(
			&p.ID, &p.Slug, &p.First, &p.Last, &p.BirthDate, &p.Description, &p.ProfileImg,
			&cm.Position, &cm.ThumbnailName, &ch.ID, &ch.Name, &ch.Image,
		)
		if err != nil {
			return nil, err
		}

		cm.Actor = p
		cm.Character = ch
		members = append(members, cm)
	}

	return members, nil
}

// TODO: make DRY later
func (m *VideoModel) GetAll(limit int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		COALESCE(c.id, 0), COALESCE(c.name, ''), COALESCE(c.page_url,''), COALESCE(c.slug,''), COALESCE(c.profile_img,'')
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		LIMIT $1;
	`
	rows, err := m.DB.Query(context.Background(), stmt, limit)
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
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
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

func (m *VideoModel) HasLike(vidId, userId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM likes WHERE video_id = $1 AND user_id = $2)"
	err := m.DB.QueryRow(context.Background(), stmt, vidId, userId).Scan(&exists)
	return exists, err
}

func (m *VideoModel) Insert(video *Video) (int, error) {
	stmt := `
	INSERT INTO video (title, video_url, upload_date, pg_rating, slug)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id;`
	result := m.DB.QueryRow(
		context.Background(), stmt, video.Title,
		video.URL, video.UploadDate, video.Rating,
		video.Slug,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *VideoModel) InsertThumbnailName(vidId int, name string) error {
	stmt := `UPDATE video SET thumbnail_name = $1 WHERE id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, name, vidId)
	return err
}

func (m *VideoModel) InsertVideoCreatorRelation(vidId, creatorId int) error {
	stmt := `INSERT INTO video_creator_rel (video_id, creator_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, creatorId)
	return err
}

func (m *VideoModel) UpdateCreatorRelation(vidId, creatorId int) error {
	stmt := `UPDATE video_creator_rel SET creator_id = $1 WHERE video_id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, creatorId)
	return err
}

func (m *VideoModel) InsertVideoPersonRelation(vidId, personId, position int, characterId *int, imgName string) error {
	stmt := `INSERT INTO video_person_rel (video_id, person_id, characterId, position, img_name) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, personId, characterId, position, imgName)
	return err
}

func (m *VideoModel) Search(query string, limit, offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, 
			v.title AS name, 
			v.slug, 
			v.thumbnail_name AS img, 
			v.upload_date, 
			c.name AS creator, 
			c.slug AS creator_slug, 
			c.profile_img AS creator_img,
			ts_rank(v.search_vector, plainto_tsquery('english', $1)) AS rank
		FROM video as v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.search_vector @@ plainto_tsquery('english', $1)
		ORDER BY rank DESC, name ASC
		LIMIT $2
		OFFSET $3;
	`
	rows, err := m.DB.Query(context.Background(), stmt, query, limit, offset)
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
			&v.ID, &v.Title, &v.Slug, &v.ThumbnailName, &v.UploadDate,
			&c.Name, &c.Slug, &c.ProfileImage, nil,
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

func (m *VideoModel) SearchCount(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM video as v
		WHERE v.search_vector @@ plainto_tsquery('english', $1)
	`
	var count int
	row := m.DB.QueryRow(context.Background(), stmt, query)
	err := row.Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	return count, nil
}

func (m *VideoModel) Update(video *Video) error {
	stmt := `
	UPDATE video SET title = $1, video_url = $2, upload_date = $3, 
	pg_rating = $4, slug = $5, thumbnail_name = $6
	WHERE id = $7`
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	result := tx.QueryRow(
		context.Background(), stmt, video.Title,
		video.URL, video.UploadDate, video.Rating,
		video.Slug, video.ThumbnailName, video.ID,
	)

	var id int
	err = result.Scan(&id)
	if err != nil {
		return err
	}

	video.ID = id
	stmt = `INSERT INTO video_creator_rel (video_id, creator_id) VALUES ($1, $2)`
	_, err = tx.Exec(context.Background(), stmt, video.ID, video.Creator.ID)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO video_person_rel (video_id, person_id, character_id, position, thumbnail) VALUES ($1, $2, $3, $4, $5)`
	for i, v := range video.Cast {
		fmt.Printf("Video:%d Person:%d Character:%d\n", video.ID, *v.Actor.ID, *v.Character.ID)
		_, err = tx.Exec(context.Background(), stmt, video.ID, v.Actor.ID, v.Character.ID, i, v.ThumbnailName)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
