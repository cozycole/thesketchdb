package models

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Show struct {
	ID         *int
	Name       *string
	Slug       *string
	ProfileImg *string
	Creator    *Creator
	Seasons    []*Season
}

type Season struct {
	ID       *int
	Number   *int
	ShowId   *int
	Episodes []*Episode
	AirDate  *time.Time
}

type Episode struct {
	ID        *int
	Number    *int
	Title     *string
	AirDate   *time.Time
	Thumbnail *string
	SeasonId  *int
	Videos    []*Video
}

type ShowModelInterface interface {
	AddSeason(showId int) (int, error)
	DeleteEpisode(episodeId int) error
	Get(filter *Filter) ([]*Show, error)
	GetById(id int) (*Show, error)
	GetBySlug(slug string) (*Show, error)
	GetEpisode(episodeId int) (*Episode, error)
	GetSeason(seasonId int) (*Season, error)
	GetShowCast(id int) ([]*Person, error)
	Insert(show *Show) (int, error)
	InsertEpisode(episode *Episode) (int, error)
	Delete(show *Show) error
	Update(show *Show) error
	UpdateEpisode(episode *Episode) error
}

type ShowModel struct {
	DB *pgxpool.Pool
}

func (m *ShowModel) AddSeason(showId int) (int, error) {
	stmt := `
		WITH latest_season AS (
			SELECT max(season_number) as last_season
			FROM season
			WHERE show_id = $1
		)
		INSERT INTO season (show_id, season_number)
		VALUES ($1, COALESCE((select last_season from latest_season), 0) + 1)
		RETURNING id
	`
	var id int
	err := m.DB.QueryRow(context.Background(), stmt, showId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShowModel) GetEpisode(episodeId int) (*Episode, error) {
	stmt := `
		SELECT e.id, e.episode_number, e.title, e.air_date, e.thumbnail_name,
		v.id, v.title, v.slug, v.video_url, v.youtube_id, v.thumbnail_name, v.upload_date, v.description, v.pg_rating,
		c.id, c.name, c.slug, c.profile_img
		FROM episode as e
		LEFT JOIN video as v ON e.id = v.episode_id
		LEFT JOIN video_creator_rel as vcr ON v.id = vcr.video_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		WHERE e.id = $1
	`

	rows, err := m.DB.Query(context.Background(), stmt, episodeId)
	if err != nil {
		return nil, err
	}

	e := &Episode{}
	for rows.Next() {
		v := &Video{}
		c := &Creator{}
		rows.Scan(
			&e.ID, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&v.ID, &v.Title, &v.Slug, &v.URL, &v.YoutubeID,
			&v.ThumbnailName, &v.UploadDate, &v.Description,
			&v.Rating, &c.ID, &c.Name, &c.Slug, &c.ProfileImage,
		)

		if v.ID == nil {
			continue
		}
		v.Creator = c

		e.Videos = append(e.Videos, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return e, nil
}

func (m *ShowModel) GetSeason(seasonId int) (*Season, error) {
	stmt := `
		SELECT DISTINCT se.id, se.season_number, se.air_date,
		e.id, e.episode_number, e.air_date, e.thumbnail_name, e.title, v.id
		FROM season as se 
		LEFT JOIN episode as e ON se.id = e.season_id
		LEFT JOIN video as v ON e.id = v.episode_id
		WHERE se.id = $1
	`

	rows, err := m.DB.Query(context.Background(), stmt, seasonId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	episodes := map[int]*Episode{}
	s := &Season{}
	for rows.Next() {
		e := &Episode{}
		v := &Video{}
		err := rows.Scan(
			&s.ID, &s.Number, &s.AirDate,
			&e.ID, &e.Number, &e.AirDate,
			&e.Thumbnail, &e.Title, &v.ID,
		)
		if err != nil {
			return nil, err
		}

		if e.ID == nil {
			continue
		}

		// If episode already exists, want to append its videos
		if currEpisode, ok := episodes[*e.ID]; ok {
			e = currEpisode
		}

		if v.ID != nil {
			e.Videos = append(e.Videos, v)
		}

		episodes[*e.ID] = e
	}
	fmt.Printf("%+v\n", episodes)
	for _, ep := range episodes {
		s.Episodes = append(s.Episodes, ep)
	}

	sort.Slice(s.Episodes, func(i, j int) bool {
		return *s.Episodes[i].Number < *s.Episodes[j].Number
	})

	return s, nil
}

func (m *ShowModel) InsertEpisode(episode *Episode) (int, error) {
	stmt := `
		INSERT INTO episode (season_id, episode_number, title, air_date, thumbnail_name)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`

	var id int
	err := m.DB.QueryRow(
		context.Background(), stmt,
		episode.SeasonId, episode.Number, episode.Title,
		episode.AirDate, episode.Thumbnail,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShowModel) Delete(show *Show) error {
	stmt := `
		DELETE FROM show
		WHERE id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, show.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *ShowModel) DeleteEpisode(episodeId int) error {
	stmt := `
		DELETE FROM episode
		WHERE id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, episodeId)
	if err != nil {
		return err
	}
	return nil
}

func (m *ShowModel) Get(filter *Filter) ([]*Show, error) {
	// stmt := `
	// 	SELECT s.id, s.name, s.profile_img, s.slug,
	// `
	return nil, nil

}

func (m *ShowModel) GetById(id int) (*Show, error) {
	stmt := `
		SELECT DISTINCT s.id, s.name, s.profile_img, s.slug,
		se.id, se.season_number, se.air_date,
		e.id, e.episode_number, e.title, e.air_date, e.thumbnail_name,
		v.id
		FROM show as s
		LEFT JOIN season as se ON s.id = se.show_id
		LEFT JOIN episode as e ON se.id = e.season_id
		LEFT JOIN video as v ON e.id = v.episode_id
		WHERE s.id = $1
	`

	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	show := &Show{}
	seasonMap := map[int]*Season{}
	seasonEpisodes := map[int]map[int]*Episode{}
	episodes := map[int]*Episode{}
	for rows.Next() {
		s := &Season{}
		e := &Episode{}
		v := &Video{}
		err := rows.Scan(
			&show.ID, &show.Name, &show.ProfileImg, &show.Slug,
			&s.ID, &s.Number, &s.AirDate,
			&e.ID, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&v.ID,
		)
		if err != nil {
			return nil, err
		}

		// no need to append anything to the struct
		// if there aren't seasons to join
		if s.ID == nil {
			continue
		}

		// add season to map if not already added
		if _, ok := seasonMap[*s.ID]; !ok {
			seasonMap[*s.ID] = s
		}

		s = seasonMap[*s.ID]

		if e.ID == nil {
			continue
		}

		if _, ok := episodes[*e.ID]; !ok {
			episodes[*e.ID] = e
		}

		e = episodes[*e.ID]

		if v.ID != nil {
			e.Videos = append(e.Videos, v)
		}

		if _, ok := seasonEpisodes[*s.ID]; !ok {
			seasonEpisodes[*s.ID] = map[int]*Episode{}
		}

		seasonEpisodes[*s.ID][*e.ID] = e
	}

	for seasonId, episodeMap := range seasonEpisodes {
		var episodes []*Episode
		for _, ep := range episodeMap {
			episodes = append(episodes, ep)
		}

		seasonMap[seasonId].Episodes = episodes
	}

	var seasons []*Season
	for _, season := range seasonMap {
		seasons = append(seasons, season)
		sort.Slice(season.Episodes, func(i, j int) bool {
			return *season.Episodes[i].Number < *season.Episodes[j].Number
		})
	}

	sort.Slice(seasons, func(i, j int) bool {
		return *seasons[i].Number < *seasons[j].Number
	})

	show.Seasons = seasons

	return show, nil
}

func (m *ShowModel) GetBySlug(slug string) (*Show, error) {
	id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.GetById(id)
}

func (m *ShowModel) GetIdBySlug(slug string) (int, error) {
	stmt := `SELECT s.id FROM show as s WHERE s.slug = $1`
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

func (m *ShowModel) GetShowCast(id int) ([]*Person, error) {
	stmt := `
		SELECT DISTINCT p.id, p.first, p.last, p.profile_img, p.birthdate, p.slug
		FROM person as p
		JOIN cast_members as cm ON p.id = cm.person_id 
		JOIN video as v ON cm.video_id = v.id
		JOIN episode as e ON v.episode_id = e.id
		JOIN season as se ON e.season_id = se.id
		JOIN show as sh ON se.show_id = sh.id
		WHERE sh.id = $1
		AND cm.role = 'cast'
		`

	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}

	var people []*Person
	for rows.Next() {
		var p Person
		destinations := []any{
			&p.ID, &p.First, &p.Last, &p.ProfileImg, &p.BirthDate, &p.Slug,
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}

		people = append(people, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return people, nil
}

func (m *ShowModel) Insert(show *Show) (int, error) {
	stmt := `
	INSERT INTO show (name, slug, profile_img)
	VALUES ($1,$2,$3)
	RETURNING id`
	result := m.DB.QueryRow(
		context.Background(), stmt, show.Name, show.Slug, show.ProfileImg,
	)
	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShowModel) Update(show *Show) error {
	stmt := `
	UPDATE show SET name = $1, slug = $2, profile_img = $3
	WHERE id = $4`
	_, err := m.DB.Exec(
		context.Background(), stmt, show.Name, show.Slug, show.ProfileImg, show.ID,
	)
	return err
}

func (m *ShowModel) UpdateEpisode(episode *Episode) error {
	stmt := `
		UPDATE episode 
		SET episode_number = $1, title = $2, air_date = $3, thumbnail_name = $4
		WHERE id = $5
	`

	_, err := m.DB.Exec(
		context.Background(), stmt, episode.Number,
		episode.Title, episode.AirDate, episode.Thumbnail,
		episode.ID,
	)
	return err
}
