package models

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"mime/multipart"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Filter struct {
	Characters []*Character
	Creators   []*Creator
	Limit      int
	Offset     int
	People     []*Person
	Query      string
	Shows      []*Show
	SortBy     string
	Tags       []*Tag
}

var sortMap = map[string]string{
	"latest": "upload_date DESC, sketch_title ASC",
	"oldest": "upload_date ASC, sketch_title ASC",
	"az":     "sketch_title ASC",
	"za":     "sketch_title DESC",
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

	for _, p := range f.Characters {
		if p.ID != nil {
			params.Add("character", strconv.Itoa(*p.ID))
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

type Sketch struct {
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
	SeasonNumber  *int
	EpisodeNumber *int
	Number        *int
	Liked         *bool
}

type SketchModelInterface interface {
	BatchUpdateTags(vidId int, tags *[]*Tag) error
	Get(filter *Filter) ([]*Sketch, error)
	GetById(id int) (*Sketch, error)
	GetCount(filter *Filter) (int, error)
	GetBySlug(slug string) (*Sketch, error)
	GetByUserLikes(id int) ([]*Sketch, error)
	GetFeatured() ([]*Sketch, error)
	HasLike(vidId, userId int) (bool, error)
	Insert(sketch *Sketch) (int, error)
	InsertThumbnailName(vidId int, name string) error
	InsertSketchCreatorRelation(vidId, creatorId int) error
	Search(search string, limit, offset int) ([]*Sketch, error)
	SearchCount(query string) (int, error)
	IsSlugDuplicate(vidId int, slug string) bool
	Update(sketch *Sketch) error
	UpdateCreatorRelation(vidId, creatorId int) error
}

type SketchModel struct {
	DB *pgxpool.Pool
}

func (m *SketchModel) BatchUpdateTags(vidId int, tags *[]*Tag) error {
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
		JOIN sketch_tags as vt
		ON t.id = vt.tag_id
		JOIN sketch as v
		ON vt.sketch_id = v.id
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
		query := "INSERT INTO sketch_tags (sketch_id, tag_id) VALUES "
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
		query := "DELETE FROM sketch_tags WHERE sketch_id = $1 AND tag_id IN ("
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

type Arguements struct {
	Args     []any
	ArgIndex int
	ImgField string
}

func determineImageField(filter *Filter) string {
	imgField := "v.thumbnail_name"
	// only use the cast thumbnail if searching for a specific person, if multiple
	// are being searched, use the overall thumbnail
	if len(filter.People) == 1 || len(filter.Characters) == 1 {
		imgField = "cm.img_name"
	}
	return imgField
}

func determineFields(filter *Filter, args *Arguements) string {
	rankParam := ""
	if filter.Query != "" {
		args.ArgIndex++
		rankParam = fmt.Sprintf(`
			 , ts_rank(
				setweight(to_tsvector('english', v.title), 'A') ||
				setweight(to_tsvector('english', c.name), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT a.first
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.sketch_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT a.last
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.sketch_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					WHERE cm.sketch_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					JOIN character AS c ON cm.character_id = c.id 
					WHERE cm.sketch_id = v.id
				), ' ')), 'B') ||
				setweight(to_tsvector('english', array_to_string(ARRAY(
					SELECT t.name
					FROM sketch_tags AS vt 
					JOIN tags AS t ON vt.tag_id = t.id 
					WHERE vt.sketch_id = v.id
				), ' ')), 'C'),
				to_tsquery('english', $%d)
				) AS rank
		`, args.ArgIndex)
		args.Args = append(args.Args, filter.Query)
	}

	baseFields := `
		v.id as sketch_id, v.title as sketch_title, v.sketch_number as sketch_number,
		v.sketch_url as sketch_url, v.slug as sketch_slug, %s as thumbnail_name, v.upload_date as upload_date, 
		c.id as creator_id, c.name as creator_name, c.slug as creator_slug, 
		c.profile_img as creator_img, sh.id as show_id, sh.name as show_name,
		sh.profile_img as show_img, sh.slug as show_slug, 
		se.season_number as season_number, e.episode_number as episode_number %s
	`

	fields := fmt.Sprintf(baseFields, args.ImgField, rankParam)

	return fields
}

func determineConditions(filter *Filter, args *Arguements) string {
	clause := ""

	if filter.Query != "" {
		args.ArgIndex++
		clause += fmt.Sprintf(`
			AND
			to_tsvector(
				'english',
				COALESCE(v.title, '') || ' ' || COALESCE(c.name, '') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT a.first
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.sketch_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT a.last
					FROM cast_members AS cm 
					JOIN person AS a ON cm.person_id = a.id 
					WHERE cm.sketch_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT t.name 
					FROM sketch_tags AS vt 
					JOIN tags AS t ON vt.tag_id = t.id 
					WHERE vt.sketch_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT cm.character_name
					FROM cast_members AS cm 
					WHERE cm.sketch_id = v.id
				), ' '),'') || ' ' ||
				COALESCE(array_to_string(ARRAY(
					SELECT c.name
					FROM cast_members AS cm 
					JOIN character as c ON cm.character_id = c.id
					WHERE cm.sketch_id = v.id
				), ' '),'')) @@ to_tsquery('english', $%d)
			`, args.ArgIndex)
		args.Args = append(args.Args, filter.Query)
	}

	// NOTE: Creators, shows and tags use OR operator
	if len(filter.Creators) > 0 {
		creatorPlaceholders := []string{}
		for _, creator := range filter.Creators {
			args.ArgIndex++
			creatorPlaceholders = append(creatorPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, creator.ID)
		}

		clause += fmt.Sprintf(" AND c.id IN (%s)", strings.Join(creatorPlaceholders, ","))
	}

	if len(filter.Characters) > 0 {
		characterPlaceholders := []string{}
		for _, character := range filter.Characters {
			args.ArgIndex++
			characterPlaceholders = append(characterPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, character.ID)
		}

		clause += fmt.Sprintf(" AND cm.character_id IN (%s)", strings.Join(characterPlaceholders, ","))
	}

	if len(filter.Shows) > 0 {
		showPlaceholders := []string{}
		for _, show := range filter.Shows {
			args.ArgIndex++
			showPlaceholders = append(showPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, show.ID)
		}

		clause += fmt.Sprintf(" AND sh.id IN (%s)", strings.Join(showPlaceholders, ","))
	}

	if len(filter.Tags) > 0 {
		tagPlaceholders := []string{}
		for _, tag := range filter.Tags {
			args.ArgIndex++
			tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, tag.ID)
		}

		clause += fmt.Sprintf(" AND vt.tag_id IN (%s)", strings.Join(tagPlaceholders, ","))
	}

	// NOTE: People filter use AND operation
	if len(filter.People) > 0 {
		peoplePlaceholders := []string{}
		for _, person := range filter.People {
			args.ArgIndex++
			peoplePlaceholders = append(peoplePlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, person.ID)
		}

		clause += fmt.Sprintf(" AND cm.person_id IN (%s)", strings.Join(peoplePlaceholders, ","))
		clause += fmt.Sprintf(`
		GROUP BY v.id, v.title, v.sketch_url, v.slug,
		         %s, v.upload_date, v.sketch_number,
		         c.id, c.name, c.page_url, c.slug, c.profile_img,
				sh.id, sh.name, sh.profile_img, sh.slug, se.season_number, e.episode_number`, args.ImgField)

		if filter.Query != "" {
			clause += ", rank"
		}

		if len(filter.People) > 1 {
			args.ArgIndex++
			clause += fmt.Sprintf(" HAVING COUNT(DISTINCT cm.person_id) = $%d ", args.ArgIndex)
			args.Args = append(args.Args, len(filter.People))
		}
	}

	return clause
}

func determineSort(filter *Filter, args *Arguements) string {
	sort := "upload_date ASC, sketch_title ASC"
	if val, ok := sortMap[filter.SortBy]; ok {
		sort = val
	}

	sort = fmt.Sprintf(" ORDER BY %s", sort)
	sort += fmt.Sprintf(" LIMIT $%d OFFSET $%d", args.ArgIndex+1, args.ArgIndex+2)
	args.ArgIndex += 2
	args.Args = append(args.Args, filter.Limit, filter.Offset)

	return sort
}

func (m *SketchModel) Get(filter *Filter) ([]*Sketch, error) {
	// The CTE is used due to possiblility of a single cast member playing
	// multiple rows in a sketch, this can cause duplicate sketch results (one for
	// each character/person pairing) so we want to limit it to one (rn = 1)
	query := `
		WITH sketch_cast AS (
		SELECT %s
		FROM sketch as v
		LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN cast_members as cm ON v.id = cm.sketch_id
		LEFT JOIN sketch_tags as vt ON v.id = vt.sketch_id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		WHERE 1=1
		%s
		),
		ranked_sketches AS (
			SELECT *, 
			ROW_NUMBER() OVER (PARTITION BY sketch_id ORDER BY sketch_id) AS rn	
			FROM sketch_cast
		)
		SELECT sketch_id, sketch_title, sketch_number, sketch_url, 
		sketch_slug, thumbnail_name, upload_date, 
		creator_id, creator_name, creator_slug, 
		creator_img, show_id, show_name,
		show_img, show_slug, season_number, episode_number %s
		FROM ranked_sketches
		WHERE rn = 1
		%s
	`

	args := &Arguements{ArgIndex: 0}

	imgField := determineImageField(filter)
	args.ImgField = imgField
	rank := ""
	if filter.Query != "" {
		rank = ", rank"
	}

	fields := determineFields(filter, args)
	conditionClause := determineConditions(filter, args)
	sortClause := determineSort(filter, args)

	fmt.Println("SORT: ", sortClause)
	query = fmt.Sprintf(query, fields, conditionClause, rank, sortClause)
	fmt.Println(query)
	fmt.Printf("ARGS: %+v\n", args.Args)

	rows, err := m.DB.Query(context.Background(), query, args.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	sketches := []*Sketch{}

	for rows.Next() {
		v := &Sketch{}
		c := &Creator{}
		sh := &Show{}
		destinations := []any{
			&v.ID, &v.Title, &v.Number, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug, &v.SeasonNumber, &v.EpisodeNumber,
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
		sketches = append(sketches, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sketches, nil
}

func (m *SketchModel) GetById(id int) (*Sketch, error) {
	stmt := `
		SELECT v.id, v.title, v.sketch_number, v.sketch_url, 
		v.slug, v.thumbnail_name, v.upload_date, v.youtube_id,
		se.season_number, e.episode_number,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		p.id, p.slug, p.first, p.last, p.profile_img,
		ch.id, ch.name, ch.slug, ch.img_name,
		cm.id, cm.position, cm.character_name, cm.img_name
		FROM sketch AS v
		LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN cast_members as cm ON v.id = cm.sketch_id
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

	v := &Sketch{}
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
			&v.ID, &v.Title, &v.Number, &v.URL, &v.Slug, &v.ThumbnailName,
			&v.UploadDate, &v.YoutubeID, &v.SeasonNumber, &v.EpisodeNumber,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&ch.ID, &ch.Name, &ch.Slug, &ch.Image,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.ThumbnailName,
		)
		if err != nil {
			return nil, err
		}

		if cm.ID != nil {
			if p.ID != nil {
				cm.Actor = p
			}

			if ch.ID != nil {
				cm.Character = ch
			}
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

func (m *SketchModel) GetBySlug(slug string) (*Sketch, error) {
	id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.GetById(id)
}

func (m *SketchModel) GetCount(filter *Filter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM (
			SELECT DISTINCT %s
			FROM sketch as v
			LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
			LEFT JOIN creator as c ON vcr.creator_id = c.id
			LEFT JOIN cast_members as cm ON v.id = cm.sketch_id
			LEFT JOIN sketch_tags as vt ON v.id = vt.sketch_id
			LEFT JOIN episode as e ON v.episode_id = e.id
			LEFT JOIN season as se ON e.season_id = se.id
			LEFT JOIN show as sh ON se.show_id = sh.id
			WHERE 1=1
			%s
		) as grouped_content
	`

	args := &Arguements{ArgIndex: 0}
	args.ImgField = "v.thumbnail_name"
	fields := determineFields(filter, args)
	conditionClause := determineConditions(filter, args)
	query = fmt.Sprintf(query, fields, conditionClause)
	fmt.Println(query)

	var count int
	err := m.DB.QueryRow(context.Background(), query, args.Args...).Scan(&count)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}

	return count, nil
}

func (m *SketchModel) GetFeatured() ([]*Sketch, error) {
	stmt := `
		SELECT v.id, v.title, v.sketch_url, v.slug, v.thumbnail_name, v.upload_date, v.youtube_id,
			c.id, c.name, c.profile_img,
			sh.id, sh.name, sh.profile_img, sh.slug,
			p.id, p.slug, p.first, p.last, p.profile_img,
			cm.id, cm.position, cm.img_name, cm.character_name,
			ch.id, ch.slug, ch.name, ch.img_name
		FROM sketch AS v
		JOIN sketch_tags as vt ON v.id = vt.sketch_id
		JOIN tags as t ON vt.tag_id = t.id
		LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN cast_members as cm ON v.id = cm.sketch_id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN character as ch ON cm.character_id = ch.id
		WHERE t.name = 'Featured'
	`

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	sketchMap := make(map[int]*Sketch)
	hasRows := false
	for rows.Next() {
		v := &Sketch{}
		c := &Creator{}
		sh := &Show{}
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		hasRows = true

		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate, &v.YoutubeID,
			&c.ID, &c.Name, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&cm.ID, &cm.Position, &cm.ThumbnailName, &cm.CharacterName,
			&ch.ID, &ch.Slug, &ch.Name, &ch.Image,
		)
		if err != nil {
			return nil, err
		}

		v.Show = sh
		v.Creator = c

		if cm.ID != nil {
			if p.ID != nil {
				cm.Actor = p
			}

			if ch.ID != nil {
				cm.Character = ch
			}
		}

		if currentVid, ok := sketchMap[*v.ID]; ok {
			currentVid.Cast = append(currentVid.Cast, cm)
		} else {
			v.Cast = append(v.Cast, cm)
			sketchMap[*v.ID] = v
		}
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	sketches := slices.Collect(maps.Values(sketchMap))

	return sketches, nil
}

func (m *SketchModel) GetIdBySlug(slug string) (int, error) {
	stmt := `SELECT v.id FROM sketch as v WHERE v.slug = $1`
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

func (m *SketchModel) GetByUserLikes(userId int) ([]*Sketch, error) {
	stmt := `
		SELECT v.id as sketch_id, v.title as sketch_title, v.sketch_number as sketch_number,
		v.sketch_url as sketch_url, v.slug as sketch_slug, v.thumbnail_name as thumbnail_name, v.upload_date as upload_date, 
		c.id as creator_id, c.name as creator_name, c.slug as creator_slug, 
		c.profile_img as creator_img, sh.id as show_id, sh.name as show_name,
		sh.profile_img as show_img, sh.slug as show_slug, 
		se.season_number as season_number, e.episode_number as episode_number
		FROM sketch AS v
		LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		JOIN likes as l ON v.id = l.sketch_id
		WHERE l.user_id = $1
		ORDER BY l.created_at desc
	`

	fmt.Println(stmt)

	rows, err := m.DB.Query(context.Background(), stmt, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	sketches := []*Sketch{}
	for rows.Next() {
		v := &Sketch{}
		c := &Creator{}
		sh := &Show{}
		destinations := []any{
			&v.ID, &v.Title, &v.Number, &v.URL, &v.Slug, &v.ThumbnailName, &v.UploadDate,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug, &v.SeasonNumber, &v.EpisodeNumber,
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		v.Show = sh
		sketches = append(sketches, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sketches, nil
}

func (m *SketchModel) HasLike(vidId, userId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM likes WHERE sketch_id = $1 AND user_id = $2)"
	err := m.DB.QueryRow(context.Background(), stmt, vidId, userId).Scan(&exists)
	return exists, err
}

func (m *SketchModel) Insert(sketch *Sketch) (int, error) {
	stmt := `
	INSERT INTO sketch (title, sketch_url, upload_date, slug)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id;`
	result := m.DB.QueryRow(
		context.Background(), stmt, sketch.Title,
		sketch.URL, sketch.UploadDate,
		sketch.Slug,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SketchModel) InsertThumbnailName(vidId int, name string) error {
	stmt := `UPDATE sketch SET thumbnail_name = $1 WHERE id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, name, vidId)
	return err
}

func (m *SketchModel) InsertSketchCreatorRelation(vidId, creatorId int) error {
	stmt := `INSERT INTO sketch_creator_rel (sketch_id, creator_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, creatorId)
	return err
}

func (m *SketchModel) UpdateCreatorRelation(vidId, creatorId int) error {
	stmt := `UPDATE sketch_creator_rel SET creator_id = $1 WHERE sketch_id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, creatorId, vidId)
	return err
}

func (m *SketchModel) Search(query string, limit, offset int) ([]*Sketch, error) {
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
		FROM sketch as v
		LEFT JOIN sketch_creator_rel as vcr
		ON v.id = vcr.sketch_id
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

	sketches := []*Sketch{}
	for rows.Next() {
		v := &Sketch{}
		c := &Creator{}
		err := rows.Scan(
			&v.ID, &v.Title, &v.Slug, &v.ThumbnailName, &v.UploadDate,
			&c.Name, &c.Slug, &c.ProfileImage, nil,
		)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		sketches = append(sketches, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sketches, nil
}

func (m *SketchModel) SearchCount(query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM sketch as v
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

func (m *SketchModel) IsSlugDuplicate(vidId int, slug string) bool {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM sketch WHERE slug = $1 AND id != $2)"
	m.DB.QueryRow(context.Background(), stmt, slug, vidId).Scan(&exists)
	return exists
}

func (m *SketchModel) Update(sketch *Sketch) error {
	stmt := `
	UPDATE sketch SET title = $1, sketch_url = $2, upload_date = $3, 
	slug = $4, thumbnail_name = $5
	WHERE id = $7`
	_, err := m.DB.Exec(
		context.Background(), stmt, sketch.Title,
		sketch.URL, sketch.UploadDate,
		sketch.Slug, sketch.ThumbnailName, sketch.ID,
	)
	return err
}
