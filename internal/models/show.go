package models

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Show struct {
	ID         *int
	Name       *string
	Aliases    *string
	Slug       *string
	ProfileImg *string
	Seasons    []*Season
}

type ShowRef struct {
	ID         *int    `json:"id"`
	Slug       *string `json:"slug"`
	Name       *string `json:"name"`
	ProfileImg *string `json:"profileImage"`
}

type Season struct {
	ID       *int
	Slug     *string
	Number   *int
	Show     *ShowRef
	Episodes []*EpisodeRef
}

type SeasonRef struct {
	ID     *int     `json:"id"`
	Slug   *string  `json:"slug"`
	Number *int     `json:"number"`
	Show   *ShowRef `json:"show"`
}

type Episode struct {
	ID        *int         `json:"id"`
	Slug      *string      `json:"slug"`
	Title     *string      `json:"title"`
	Number    *int         `json:"number"`
	AirDate   *time.Time   `json:"airDate"`
	Thumbnail *string      `json:"thumbnail"`
	Season    *SeasonRef   `json:"season"`
	URL       *string      `json:"url"`
	Sketches  []*SketchRef `json:"sketches"`
	YoutubeID *string      `json:"youtubeId"`
}

func (e *Episode) GetTitle() *string     { return e.Title }
func (e *Episode) GetNumber() *int       { return e.Number }
func (e *Episode) GetSeason() *SeasonRef { return e.Season }
func (e *Episode) GetSketchCount() int   { return len(e.Sketches) }

func (e *Episode) GetShow() *ShowRef {
	if e.Season == nil ||
		e.Season.ID == nil ||
		e.Season.Show == nil ||
		e.Season.Show.ID == nil {
		return nil
	}
	return e.Season.Show
}

type EpisodeRef struct {
	ID          *int       `json:"id"`
	Slug        *string    `json:"slug"`
	Title       *string    `json:"title"`
	Number      *int       `json:"number"`
	AirDate     *time.Time `json:"airDate"`
	Thumbnail   *string    `json:"thumbnail"`
	Season      *SeasonRef `json:"season"`
	SketchCount *int       `json:"-"`
}

func (e *EpisodeRef) GetTitle() *string     { return e.Title }
func (e *EpisodeRef) GetNumber() *int       { return e.Number }
func (e *EpisodeRef) GetSeason() *SeasonRef { return e.Season }
func (e *EpisodeRef) GetSketchCount() int   { return safeDeref(e.SketchCount) }

func (e *EpisodeRef) GetShow() *ShowRef {
	if e.Season == nil ||
		e.Season.ID == nil ||
		e.Season.Show.ID == nil {
		return nil
	}
	return e.Season.Show
}

func (s *Season) AirYear() string {
	var airDates []time.Time
	for _, e := range s.Episodes {
		if e.AirDate != nil {
			airDates = append(airDates, *e.AirDate)
		}
	}
	if len(airDates) == 0 {
		return ""
	}

	min := airDates[0]
	for _, t := range airDates[1:] {
		if t.Before(min) {
			min = t
		}
	}
	return strconv.Itoa(min.Year())
}

type ShowModelInterface interface {
	AddSeason(*Season) (int, error)
	Delete(show *Show) error
	DeleteEpisode(episodeId int) error
	DeleteSeason(seasonid int) error
	EpisodeExists(id int) (bool, error)
	Get(filter *Filter) ([]*Show, error)
	GetById(id int) (*Show, error)
	GetBySlug(slug string) (*Show, error)
	GetCount(filter *Filter) (int, error)
	GetEpisode(episodeId int) (*Episode, error)
	GetSeason(seasonId int) (*Season, error)
	GetShowCast(id int) ([]*Person, error)
	GetShowRefs(ids []int) ([]*ShowRef, error)
	Insert(show *Show) (int, error)
	InsertEpisode(episode *Episode) (int, error)
	ListEpisodes(f *Filter) ([]*EpisodeRef, Metadata, error)
	Search(query string) ([]*Show, error)
	Update(show *Show) error
	UpdateEpisode(episode *Episode) error
}

type ShowModel struct {
	DB *pgxpool.Pool
}

func (m *ShowModel) AddSeason(season *Season) (int, error) {
	stmt := `
		INSERT INTO season (show_id, season_number, slug)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := m.DB.QueryRow(context.Background(), stmt, season.Show.ID, season.Number, season.Slug).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShowModel) GetEpisode(episodeId int) (*Episode, error) {
	stmt := `
		SELECT e.id, e.slug, e.episode_number, e.title, e.air_date, 
		e.thumbnail_name, e.url, e.youtube_id,
		v.id, v.title, v.slug, v.sketch_number, 
		v.thumbnail_name, v.upload_date, v.rating,
		se.id, se.slug, se.season_number,
		sh.id, sh.name, sh.profile_img, sh.slug
		FROM episode as e
		LEFT JOIN sketch as v ON e.id = v.episode_id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		WHERE e.id = $1
		ORDER BY v.sketch_number asc
	`

	rows, err := m.DB.Query(context.Background(), stmt, episodeId)
	if err != nil {
		return nil, err
	}

	e := &Episode{}
	se := &SeasonRef{}
	sh := &ShowRef{}
	for rows.Next() {
		v := &SketchRef{}
		rows.Scan(
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&e.URL, &e.YoutubeID,
			&v.ID, &v.Title, &v.Slug, &v.Number,
			&v.Thumbnail, &v.UploadDate, &v.Rating,
			&se.ID, &se.Slug, &se.Number,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
		)

		if v.ID == nil {
			continue
		}
		v.Episode = &EpisodeRef{
			ID:        e.ID,
			Slug:      e.Slug,
			Title:     e.Title,
			Number:    e.Number,
			AirDate:   e.AirDate,
			Thumbnail: e.Thumbnail,
			Season:    se,
		}

		e.Sketches = append(e.Sketches, v)
	}

	se.Show = sh
	e.Season = se

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return e, nil
}

func (m *ShowModel) EpisodeExists(id int) (bool, error) {
	stmt := `SELECT EXISTS(
		SELECT 1 FROM episode WHERE id = $1
	)`
	var exists bool
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&exists)
	if err == pgx.ErrNoRows {
		err = nil
	}

	return exists, err
}

func (m *ShowModel) GetSeason(seasonId int) (*Season, error) {
	stmt := `
		SELECT DISTINCT se.id, se.slug, se.season_number,
		sh.id, sh.slug, sh.name, sh.profile_img,
		e.id, e.slug, e.episode_number, e.air_date, e.thumbnail_name, e.title, v.id
		FROM season as se
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN episode as e ON se.id = e.season_id
		LEFT JOIN sketch as v ON e.id = v.episode_id
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

	episodes := map[int]*EpisodeRef{}
	s := &Season{}
	sh := &ShowRef{}
	for rows.Next() {
		e := &EpisodeRef{}
		v := &SketchRef{}
		err := rows.Scan(
			&s.ID, &s.Slug, &s.Number,
			&sh.ID, &sh.Slug, &sh.Name, &sh.ProfileImg,
			&e.ID, &e.Slug, &e.Number, &e.AirDate,
			&e.Thumbnail, &e.Title, &v.ID,
		)
		if err != nil {
			return nil, err
		}

		if e.ID == nil {
			continue
		}

		e.Season = &SeasonRef{
			ID:     s.ID,
			Slug:   s.Slug,
			Number: s.Number,
			Show:   sh,
		}

		// If episode already exists, want to append its sketches
		if currEpisode, ok := episodes[*e.ID]; ok {
			e = currEpisode
		}

		if v.ID != nil {
			newTotal := safeDeref(e.SketchCount) + 1
			e.SketchCount = &newTotal
		}

		episodes[*e.ID] = e
	}

	s.Show = sh

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
		INSERT INTO episode 
		(season_id, episode_number, title, url, air_date, thumbnail_name, youtube_id, slug)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
	`

	var id int
	err := m.DB.QueryRow(
		context.Background(), stmt,
		episode.Season.ID, episode.Number, episode.Title,
		episode.URL, episode.AirDate, episode.Thumbnail, &episode.YoutubeID, &episode.Slug,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShowModel) Delete(show *Show) error {
	if show.ID == nil {
		return fmt.Errorf("No show ID specified to delete")
	}
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

func (m *ShowModel) DeleteSeason(seasonId int) error {
	stmt := `
		DELETE FROM season
		WHERE id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, seasonId)
	if err != nil {
		return err
	}
	return nil
}

func (m *ShowModel) Get(filter *Filter) ([]*Show, error) {
	query := `
		SELECT s.id, s.slug, s.aliases, s.name, s.profile_img %s
		FROM show as s
		WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {
		rankParam := fmt.Sprintf(`
		, ts_rank(
			setweight(to_tsvector('english', s.name || 
			' ' || COALESCE(s.aliases, '')), 'A'),
			websearch_to_tsquery('english', $%d)
		) AS rank
		`, argIndex)

		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`AND
            to_tsvector('english', s.name || ' ' || COALESCE(s.aliases, '')) @@ websearch_to_tsquery('english', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := m.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	var shows []*Show
	for rows.Next() {
		var s Show
		destinations := []any{
			&s.ID, &s.Slug, &s.Aliases, &s.Name, &s.ProfileImg,
		}

		var rank *float32
		if filter.Query != "" {
			destinations = append(destinations, &rank)
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}

		shows = append(shows, &s)
	}

	return shows, nil
}

func (m *ShowModel) GetCount(filter *Filter) (int, error) {
	query := `
			SELECT COUNT(*)
			FROM (
				SELECT DISTINCT s.id, s.slug, s.name, s.profile_img
				FROM show as s
				WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if filter.Query != "" {

		query += fmt.Sprintf(`AND
            to_tsvector('english', s.name || ' ' 
			|| COALESCE(s.aliases, '')) @@ websearch_to_tsquery('english', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	query += " ) as grouped_count"

	var count int
	err := m.DB.QueryRow(context.Background(), query, args...).Scan(&count)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}

	return count, nil
}

func (m *ShowModel) GetById(id int) (*Show, error) {
	stmt := `
		SELECT DISTINCT s.id, s.name, s.aliases, s.profile_img, s.slug,
		se.id, se.slug, se.season_number, 
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name,
		v.id
		FROM show as s
		LEFT JOIN season as se ON s.id = se.show_id
		LEFT JOIN episode as e ON se.id = e.season_id
		LEFT JOIN sketch as v ON e.id = v.episode_id
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
	seasonEpisodes := map[int]map[int]*EpisodeRef{}
	episodes := map[int]*EpisodeRef{}
	for rows.Next() {
		s := &Season{}
		e := &EpisodeRef{}
		v := &SketchRef{}
		err := rows.Scan(
			&show.ID, &show.Name, &show.Aliases, &show.ProfileImg, &show.Slug,
			&s.ID, &s.Slug, &s.Number,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
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

		showRef := &ShowRef{
			ID:         show.ID,
			Slug:       show.Slug,
			Name:       show.Name,
			ProfileImg: show.ProfileImg,
		}

		s.Show = showRef
		// add season to map if not already added
		if _, ok := seasonMap[*s.ID]; !ok {
			seasonMap[*s.ID] = s
		}

		s = seasonMap[*s.ID]

		if e.ID == nil {
			continue
		}

		e.Season = &SeasonRef{
			ID:     s.ID,
			Slug:   s.Slug,
			Number: s.Number,
			Show:   showRef,
		}
		if _, ok := episodes[*e.ID]; !ok {
			episodes[*e.ID] = e
		}

		e = episodes[*e.ID]

		if v.ID != nil {
			newCount := safeDeref(e.SketchCount) + 1
			e.SketchCount = &newCount
		}

		if _, ok := seasonEpisodes[*s.ID]; !ok {
			seasonEpisodes[*s.ID] = map[int]*EpisodeRef{}
		}

		seasonEpisodes[*s.ID][*e.ID] = e
	}

	if show.ID == nil {
		return nil, ErrNoRecord
	}

	for seasonId, episodeMap := range seasonEpisodes {
		var episodes []*EpisodeRef
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
		JOIN sketch as v ON cm.sketch_id = v.id
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
	INSERT INTO show (name, aliases, slug, profile_img)
	VALUES ($1,$2,$3,$4)
	RETURNING id`
	result := m.DB.QueryRow(
		context.Background(), stmt, show.Name,
		show.Aliases, show.Slug, show.ProfileImg,
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
	UPDATE show SET name = $1, aliases = $2, slug = $3, profile_img = $4
	WHERE id = $5`
	_, err := m.DB.Exec(
		context.Background(), stmt, show.Name, show.Aliases,
		show.Slug, show.ProfileImg, show.ID,
	)
	return err
}

func (m *ShowModel) GetShowRefs(ids []int) ([]*ShowRef, error) {
	if len(ids) < 1 {
		return nil, nil
	}

	stmt := `SELECT id, name, slug, profile_img
			FROM show
			WHERE id IN (%s)`

	args := []any{}
	queryPlaceholders := []string{}
	for i, id := range ids {
		queryPlaceholders = append(queryPlaceholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	stmt = fmt.Sprintf(stmt, strings.Join(queryPlaceholders, ","))
	rows, err := m.DB.Query(context.Background(), stmt, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var shows []*ShowRef
	for rows.Next() {
		s := ShowRef{}
		err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.ProfileImg)
		if err != nil {
			return nil, err
		}
		shows = append(shows, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shows, nil
}

func (m *ShowModel) Search(query string) ([]*Show, error) {
	query = "%" + query + "%"
	stmt := `SELECT s.id, s.slug, s.name, s.profile_img
			FROM show as s
			WHERE name ILIKE $1
			ORDER BY name`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shows := []*Show{}
	for rows.Next() {
		c := &Show{}
		err := rows.Scan(
			&c.ID, &c.Slug, &c.Name, &c.ProfileImg,
		)
		if err != nil {
			return nil, err
		}
		shows = append(shows, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shows, nil
}

type EpisodeQuery struct {
	ShowName      string
	SeasonNumber  *int
	EpisodeNumber *int
}

func (m *ShowModel) ListEpisodes(f *Filter) ([]*EpisodeRef, Metadata, error) {
	stmt := `
		SELECT count(*) OVER(), e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name,
		s.id, s.slug, s.season_number,
		sh.id, sh.slug, sh.name, sh.profile_img
		FROM episode as e
		JOIN season as s ON e.season_id = s.id
		JOIN show as sh ON s.show_id = sh.id
		WHERE 
		  LOWER(sh.name) ILIKE '%' || LOWER($1) || '%'
		AND (COALESCE($2, s.season_number) = s.season_number)
		AND (COALESCE($3, e.episode_number) = e.episode_number)
		ORDER BY sh.name, s.season_number, e.episode_number
		LIMIT $4
		OFFSET $5;
	`
	epQuery, err := ExtractEpisodeQuery(f.Query)
	if err != nil {
		return nil, Metadata{}, err
	}

	fmt.Printf("%+v\n", epQuery)
	rows, err := m.DB.Query(context.Background(), stmt, epQuery.ShowName,
		epQuery.SeasonNumber, epQuery.EpisodeNumber, f.Limit(), f.Offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	episodes := []*EpisodeRef{}
	var totalCount int
	for rows.Next() {
		e := &EpisodeRef{}
		sh := &ShowRef{}
		se := &SeasonRef{}
		err := rows.Scan(
			&totalCount, &e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&se.ID, &se.Slug, &se.Number, &sh.ID, &sh.Slug, &sh.Name,
			&sh.ProfileImg,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		se.Show = sh
		e.Season = se
		episodes = append(episodes, e)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return episodes, calculateMetadata(totalCount, f.Page, f.PageSize), nil
}

func ExtractEpisodeQuery(input string) (EpisodeQuery, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))
	episodeQuery := EpisodeQuery{}

	// First check for s01e02 or s1e2 pattern
	seRe := regexp.MustCompile(`s(\d{1,2})e(\d{1,2})`)
	seMatch := seRe.FindStringSubmatch(normalized)

	if len(seMatch) == 3 {
		season, err := strconv.Atoi(seMatch[1])
		if err != nil {
			return EpisodeQuery{}, err
		}
		episode, err := strconv.Atoi(seMatch[2])
		if err != nil {
			return EpisodeQuery{}, err
		}
		episodeQuery.SeasonNumber = &season
		episodeQuery.EpisodeNumber = &episode

		normalized = strings.Replace(normalized, seMatch[0], "", 1)
	} else {
		// Fallback: check for just s01 or s1 pattern
		sRe := regexp.MustCompile(`s(\d{1,2})`)
		sMatch := sRe.FindStringSubmatch(normalized)
		if len(sMatch) == 2 {
			season, err := strconv.Atoi(sMatch[1])
			if err != nil {
				return EpisodeQuery{}, err
			}
			episodeQuery.SeasonNumber = &season

			normalized = strings.Replace(normalized, sMatch[0], "", 1)
		}
	}

	episodeQuery.ShowName = strings.TrimSpace(normalized)

	return episodeQuery, nil
}

func (m *ShowModel) UpdateEpisode(episode *Episode) error {
	stmt := `
		UPDATE episode 
		SET episode_number = $1, title = $2, air_date = $3, thumbnail_name = $4,
		url = $5, youtube_id = $6
		WHERE id = $7
	`

	_, err := m.DB.Exec(
		context.Background(), stmt, episode.Number,
		episode.Title, episode.AirDate, episode.Thumbnail,
		episode.URL, episode.YoutubeID, episode.ID,
	)
	return err
}
