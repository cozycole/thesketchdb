package models

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Filter struct {
	Query    string
	People   []*Person
	Creators []*Creator
	Shows    []*Show
	Tags     []*Tag
	SortBy   string
	Limit    int
	Offset   int
}

var sortMap = map[string]string{
	"latest": "v.upload_date DESC, v.title ASC",
	"oldest": "v.upload_date ASC, v.title ASC",
	"az":     "v.title ASC",
	"za":     "v.title DESC",
}

func (f *Filter) Params() url.Values {
	params := url.Values{}

	if f.SortBy != "" {
		params.Add("sort", f.SortBy)
	}

	if f.Query != "" {
		params.Add("query", url.QueryEscape(f.Query))
	}

	for _, p := range f.People {
		if p.ID != nil {
			params.Add("person", strconv.Itoa(*p.ID))
		}
	}

	for _, p := range f.Creators {
		if p.ID != nil {
			params.Add("creator", strconv.Itoa(*p.ID))
		}
	}

	for _, p := range f.Shows {
		if p.ID != nil {
			params.Add("show", strconv.Itoa(*p.ID))
		}
	}

	for _, p := range f.Tags {
		if p.ID != nil {
			params.Add("tag", strconv.Itoa(*p.ID))
		}
	}

	return params
}

func (f *Filter) ParamsString() string {
	return f.Params().Encode()
}

type Video struct {
	ID            *int
	Title         *string
	URL           *string
	YoutubeID     *string
	Slug          *string
	ThumbnailName *string
	ThumbnailFile *multipart.FileHeader
	Rating        *string
	Description   *string
	UploadDate    *time.Time
	Creator       *Creator
	Cast          []*CastMember
	Tags          *[]*Tag
	Show          *Show
	Number        *int
	Liked         bool
}

type VideoModelInterface interface {
	BatchUpdateTags(vidId int, tags *[]*Tag) error
	Get(filter *Filter) ([]*Video, error)
	GetById(id int) (*Video, error)
	GetCount(filter *Filter) (int, error)
	GetBySlug(slug string) (*Video, error)
	GetByUserLikes(id int) ([]*Video, error)
	HasLike(vidId, userId int) (bool, error)
	Insert(video *Video) (int, error)
	InsertThumbnailName(vidId int, name string) error
	InsertVideoCreatorRelation(vidId, creatorId int) error
	Search(search string, limit, offset int) ([]*Video, error)
	SearchCount(query string) (int, error)
	IsSlugDuplicate(vidId int, slug string) bool
	Update(video *Video) error
	UpdateCreatorRelation(vidId, creatorId int) error
}

type VideoModel struct {
	DB *pgxpool.Pool
}

func (m *VideoModel) BatchUpdateTags(vidId int, tags *[]*Tag) error {
	if tags == nil {
		return fmt.Errorf("tags argument is nil")
	}

	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Get existing tags
	stmt := `
		SELECT t.id
		FROM tags as t
		JOIN video_tags as vt
		ON t.id = vt.tag_id
		JOIN video as v
		ON vt.video_id = v.id
		WHERE v.id = $1
	`

	var id int
	existingTags := make(map[int]bool)
	rows, err := tx.Query(context.Background(), stmt, vidId)
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		existingTags[id] = true
	}

	fmt.Printf("EXISTING TAGS: %+v\n", existingTags)

	// New tags
	newTags := make(map[int]bool)
	for _, tag := range *tags {
		newTags[*tag.ID] = true
	}

	tagsToInsert := []int{}
	for tag_id := range newTags {
		if !existingTags[tag_id] {
			tagsToInsert = append(tagsToInsert, tag_id)
		}
	}
	fmt.Printf("TAGS TO INSERT: %+v\n", tagsToInsert)

	// Find tags to delete
	tagsToDelete := []int{}
	for tag_id := range existingTags {
		if !newTags[tag_id] {
			tagsToDelete = append(tagsToDelete, tag_id)
		}
	}

	if len(tagsToInsert) > 0 {
		query := "INSERT INTO video_tags (video_id, tag_id) VALUES "
		values := []interface{}{}
		for i, tag := range tagsToInsert {
			query += fmt.Sprintf("($1, $d),", i+2)
			values = append(values, tag)
		}
		query = query[:len(query)-1] // Trim last comma
		values = append([]interface{}{vidId}, values...)
		fmt.Printf("QUERY: %s\n", query)
		fmt.Printf("VALUES: %+v\n", values)

		_, err = tx.Exec(context.Background(), query, values...)
		if err != nil {
			return err
		}
	}

	if len(tagsToDelete) > 0 {
		query := "DELETE FROM video_tags WHERE video_id = $1 AND tag_id IN ("
		values := []interface{}{vidId}
		for i, tag := range tagsToDelete {
			query += fmt.Sprintf("$%d,", i+2)
			values = append(values, tag)
		}
		query = query[:len(query)-1] + ")"
		fmt.Printf("QUERY: %s", query)
		_, err = tx.Exec(context.Background(), query, values...)
		if err != nil {
			return err
		}
	}

	tx.Commit(context.Background())
	return nil
}

func (m *VideoModel) Get(filter *Filter) ([]*Video, error) {
	query := `
		SELECT DISTINCT 
		v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date, 
		c.id, c.name, c.page_url, c.slug, c.profile_img,
		sh.id, sh.name, sh.profile_img, sh.slug%s
		FROM video as v
		LEFT JOIN video_creator_rel as vcr ON v.id = vcr.video_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN cast_members as cm ON v.id = cm.video_id
		LEFT JOIN video_tags as vt ON v.id = vt.video_id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filter.Query != "" {
		rankParam := fmt.Sprintf(`
			 , ts_rank(
				setweight(to_tsvector('english', v.title), 'A') ||
				setweight(to_tsvector('english', c.name), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT a.first
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.video_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT a.last
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.video_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					WHERE cm.video_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					JOIN character AS c ON cm.character_id = c.id 
					WHERE cm.video_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT t.name
					FROM video_tags AS vt 
					JOIN tags AS t ON vt.tag_id = t.id 
					WHERE vt.video_id = v.id
				), ' ')), 'C'),
				to_tsquery('english', $%d)
				) AS rank
			`, argIndex)
		query = fmt.Sprintf(query, rankParam)

		query += fmt.Sprintf(`
			AND
            to_tsvector(
				'english',
				COALESCE(v.title, '') || ' ' || COALESCE(c.name, '') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT a.first
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT a.last
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT t.name 
					FROM video_tags AS vt 
					JOIN tags AS t ON vt.tag_id = t.id 
					WHERE vt.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					WHERE cm.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT c.name
					FROM cast_members AS cm 
					JOIN character as c ON cm.character_id = c.id
					WHERE cm.video_id = v.id
				), ' '),'')) @@ to_tsquery('english', $%d)
		
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	} else {
		query = fmt.Sprintf(query, "")
	}

	// NOTE: Creators, shows, tags filters use OR opeartions
	if len(filter.Creators) > 0 {
		creatorPlaceholders := []string{}
		for _, creator := range filter.Creators {
			creatorPlaceholders = append(creatorPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, creator.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND c.id IN (%s)", strings.Join(creatorPlaceholders, ","))
	}

	if len(filter.Shows) > 0 {
		showPlaceholders := []string{}
		for _, show := range filter.Shows {
			showPlaceholders = append(showPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, show.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND sh.id IN (%s)", strings.Join(showPlaceholders, ","))
	}

	if len(filter.Tags) > 0 {
		tagPlaceholders := []string{}
		for _, tag := range filter.Tags {
			tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, tag.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND vt.tag_id IN (%s)", strings.Join(tagPlaceholders, ","))
	}

	// NOTE: People filter use AND operation
	if len(filter.People) > 0 {
		peoplePlaceholders := []string{}
		for _, person := range filter.People {
			peoplePlaceholders = append(peoplePlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, person.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND cm.person_id IN (%s)", strings.Join(peoplePlaceholders, ","))
		query += `
		GROUP BY v.id, v.title, v.video_url, v.slug, 
		         v.thumbnail_name, v.upload_date, 
		         c.id, c.name, c.page_url, c.slug, c.profile_img,
				sh.id, sh.name, sh.profile_img, sh.slug
		`
		if len(filter.People) > 1 {
			query += fmt.Sprintf("HAVING COUNT(DISTINCT cm.person_id) = $%d ", argIndex)
			args = append(args, len(filter.People))
			argIndex++
		}
	}

	sort := "v.upload_date ASC, v.title ASC"
	if val, ok := sortMap[filter.SortBy]; ok {
		sort = val
	}

	query += fmt.Sprintf(" ORDER BY %s", sort)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)
	fmt.Println(query)
	fmt.Printf("ARGS: %+v\n", args)

	rows, err := m.DB.Query(context.Background(), query, args...)
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
		sh := &Show{}
		destinations := []any{
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
			&c.ID, &c.Name, &c.URL, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
		}
		var rank *float32
		if filter.Query != "" {
			destinations = append(destinations, &rank)
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		v.Show = sh
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) GetById(id int) (*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
			c.id, c.name, c.profile_img,
			sh.id, sh.name, sh.profile_img, sh.slug,
			p.id, p.slug, p.first, p.last, p.birthdate,
			p.description, p.profile_img,
			cm.id, cm.position, cm.img_name, cm.character_name,
			ch.id, ch.slug, ch.name, ch.img_name
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr ON v.id = vcr.video_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN cast_members as cm ON v.id = cm.video_id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN character as ch ON cm.character_id = ch.id
		WHERE v.id = $1
	`

	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	v := &Video{}
	c := &Creator{}
	sh := &Show{}
	members := []*CastMember{}
	hasRows := false
	for rows.Next() {
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		hasRows = true
		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
			&c.ID, &c.Name, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.BirthDate, &p.Description, &p.ProfileImg,
			&cm.ID, &cm.Position, &cm.ThumbnailName, &cm.CharacterName,
			&ch.ID, &ch.Slug, &ch.Name, &ch.Image,
		)
		if err != nil {
			return nil, err
		}
		if cm.ID != nil {
			cm.Actor = p
			cm.Character = ch
			members = append(members, cm)
		}
	}

	if !hasRows {
		return nil, ErrNoRecord
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	v.Show = sh
	v.Creator = c
	v.Cast = members
	return v, nil
}

func (m *VideoModel) GetBySlug(slug string) (*Video, error) {
	id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.GetById(id)
}

func (m *VideoModel) GetCount(filter *Filter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM (
			SELECT DISTINCT v.id, v.title, v.video_url, v.slug, 
			v.thumbnail_name, v.upload_date, c.id, c.name, c.page_url, 
			c.slug, c.profile_img
			FROM video as v
			LEFT JOIN video_creator_rel as vcr ON v.id = vcr.video_id
			LEFT JOIN creator as c ON vcr.creator_id = c.id
			LEFT JOIN cast_members as cm ON v.id = cm.video_id
			LEFT JOIN video_tags as vt ON v.id = vt.video_id
			WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if len(filter.Creators) > 0 {
		creatorPlaceholders := []string{}
		for _, creator := range filter.Creators {
			creatorPlaceholders = append(creatorPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, creator.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND c.id IN (%s)", strings.Join(creatorPlaceholders, ","))
	}

	if len(filter.Shows) > 0 {
		showPlaceholders := []string{}
		for _, show := range filter.Shows {
			showPlaceholders = append(showPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, show.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND sh.id IN (%s)", strings.Join(showPlaceholders, ","))
	}

	if len(filter.Tags) > 0 {
		tagPlaceholders := []string{}
		for _, tag := range filter.Tags {
			tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, tag.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND vt.tag_id IN (%s)", strings.Join(tagPlaceholders, ","))
	}

	if filter.Query != "" {
		query += fmt.Sprintf(`
			AND
            to_tsvector(
				'english',
				COALESCE(v.title, '') || ' ' || COALESCE(c.name, '') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT a.first
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT a.last
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT t.name 
					FROM video_tags AS vt 
					JOIN tags AS t ON vt.tag_id = t.id 
					WHERE vt.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					WHERE cm.video_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT c.name
					FROM cast_members AS cm 
					JOIN character as c ON cm.character_id = c.id
					WHERE cm.video_id = v.id
				), ' '),'')) @@ to_tsquery('english', $%d)
		`, argIndex)

		args = append(args, filter.Query)
		argIndex++
	}

	if len(filter.People) > 0 {
		peoplePlaceholders := []string{}
		for _, person := range filter.People {
			peoplePlaceholders = append(peoplePlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, person.ID)
			argIndex++
		}

		query += fmt.Sprintf(" AND cm.person_id IN (%s)", strings.Join(peoplePlaceholders, ","))
		query += `
		GROUP BY v.id, v.title, v.video_url, v.slug, 
		         v.thumbnail_name, v.upload_date, 
		         c.id, c.name, c.page_url, c.slug, c.profile_img
		`
		if len(filter.People) > 1 {
			query += fmt.Sprintf("HAVING COUNT(DISTINCT cm.person_id) = $%d ", argIndex)
			args = append(args, len(filter.People))
			argIndex++
		}
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
	_, err := m.DB.Exec(context.Background(), stmt, creatorId, vidId)
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
			ts_rank(v.search_vector, websearch_to_tsquery('english', $1)) AS rank
		FROM video as v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.search_vector @@ websearch_to_tsquery('english', $1)
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
		WHERE v.search_vector @@ websearch_to_tsquery('english', $1)
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

func (m *VideoModel) IsSlugDuplicate(vidId int, slug string) bool {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM video WHERE slug = $1 AND id != $2)"
	m.DB.QueryRow(context.Background(), stmt, slug, vidId).Scan(&exists)
	return exists
}

func (m *VideoModel) Update(video *Video) error {
	stmt := `
	UPDATE video SET title = $1, video_url = $2, upload_date = $3, 
	pg_rating = $4, slug = $5, thumbnail_name = $6
	WHERE id = $7`
	_, err := m.DB.Exec(
		context.Background(), stmt, video.Title,
		video.URL, video.UploadDate, video.Rating,
		video.Slug, video.ThumbnailName, video.ID,
	)
	return err
}
