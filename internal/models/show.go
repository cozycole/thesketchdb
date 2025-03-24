package models

import (
	"context"
	"errors"
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
	Season   *Show
	Episodes []*Episode
	AirDate  *time.Time
}

type Episode struct {
	ID        *int
	Number    *int
	AirDate   *time.Time
	Thumbnail *string
	Videos    []*Video
}

type ShowModelInterface interface {
	Get(filter *Filter) ([]*Show, error)
	GetById(id int) (*Show, error)
	GetBySlug(slug string) (*Show, error)
	GetShowCast(id int) ([]*Person, error)
}

type ShowModel struct {
	DB *pgxpool.Pool
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
		c.id, c.name, c.profile_img, c.slug,
		se.id, se.season_number, se.air_date,
		e.id, e.episode_number, v.id
		FROM show as s
		LEFT JOIN show_creator as sc ON s.id = sc.show_id
		LEFT JOIN creator as c on sc.creator_id = c.id
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
	c := &Creator{}
	seasonMap := map[int]*Season{}
	seasonEpisodes := map[int]map[int]*Episode{}
	episodes := map[int]*Episode{}
	for rows.Next() {
		s := &Season{}
		e := &Episode{}
		v := &Video{}
		err := rows.Scan(
			&show.ID, &show.Name, &show.ProfileImg, &show.Slug,
			&c.ID, &c.Name, &c.ProfileImage, &c.Slug,
			&s.ID, &s.Number, &s.AirDate,
			&e.ID, &e.Number, &v.ID,
		)
		if err != nil {
			return nil, err
		}

		// no need to update anything if there aren't seasons
		// to join
		if s.ID == nil {
			continue
		}

		seasonMap[*s.ID] = s

		if e.ID == nil {
			continue
		}

		// If episode already exists, want to append its videos
		if currEpisode, ok := episodes[*s.ID]; ok {
			e = currEpisode
		}

		if v.ID != nil {
			e.Videos = append(e.Videos, v)
		}

		if _, ok := seasonEpisodes[*s.ID]; !ok {
			seasonEpisodes[*s.ID] = map[int]*Episode{}
		}

		seasonEpisodes[*s.ID][*e.ID] = e
	}

	show.Creator = c

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
