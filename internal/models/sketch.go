package models

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Sketch struct {
	ID            *int          `json:"id"`
	Slug          *string       `json:"slug"`
	Title         *string       `json:"title"`
	URL           *string       `json:"url"`
	Duration      *int          `json:"duration"`
	Description   *string       `json:"description"`
	YoutubeID     *string       `json:"youtubeId"`
	ThumbnailName *string       `json:"thumbnailName"`
	Popularity    *float32      `json:"popularity"`
	UploadDate    *time.Time    `json:"uploadDate"`
	Creator       *CreatorRef   `json:"creator"`
	Cast          []*CastMember `json:"cast"`
	CastThumbnail *string       `json:"castThumbnail"`
	Tags          *[]*Tag       `json:"tags"`
	Episode       *Episode      `json:"episode"`
	EpisodeStart  *int          `json:"episodeStart"`
	Number        *int          `json:"episodeSketchOrder"`
	Series        *SeriesRef    `json:"series"`
	SeriesPart    *int          `json:"seriesPart"`
	Recurring     *RecurringRef `json:"recurring"`
	Rating        *float32      `json:"rating"`
	TotalRatings  *int          `json:"totalRatings"`
	Liked         *bool         `json:"liked,omitempty"`
}

type SketchRef struct {
	ID            *int        `json:"id"`
	Slug          *string     `json:"slug"`
	Title         *string     `json:"title"`
	Thumbnail     *string     `json:"thumbnail"`
	CastThumbnail *string     `json:"castThumbnail"`
	UploadDate    *time.Time  `json:"uploadDate"`
	Creator       *CreatorRef `json:"creator"`
	Episode       *EpisodeRef `json:"episode"`
	Number        *int        `json:"episodeSketchOrder"`
	Rating        *float32    `json:"rating"`
}

type SketchModelInterface interface {
	BatchUpdateTags(sketchId int, tags *[]*Tag) error
	Delete(id int) error
	Exists(id int) (bool, error)
	Get(filter *Filter) ([]*SketchRef, Metadata, error)
	GetById(id int) (*Sketch, error)
	GetByUserLikes(id int) ([]*SketchRef, error)
	GetCount(filter *Filter) (int, error)
	GetFeatured() ([]*Sketch, error)
	HasLike(sketchId, userId int) (bool, error)
	Insert(sketch *Sketch) (int, error)
	InsertSketchCreatorRelation(sketchId, creatorId int) error
	InsertThumbnailName(sketchId int, name string) error
	SearchCount(query string) (int, error)
	SyncSketchCreators(sketchID int, creatorIDs []int) error
	Update(sketch *Sketch) error
	UpdateCreatorRelation(sketchId, creatorId int) error
}

type SketchModel struct {
	DB *pgxpool.Pool
}

func (m *SketchModel) Delete(id int) error {
	stmt := `
		DELETE from sketch
		WHERE id = $1
	`
	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *SketchModel) Exists(id int) (bool, error) {
	stmt := `
	SELECT EXISTS (
	  SELECT 1 FROM sketch WHERE id = $1  
	);
	`

	var exists bool

	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&exists)
	return exists, err

}

func (m *SketchModel) BatchUpdateTags(sketchId int, tags *[]*Tag) error {
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
	rows, err := tx.Query(context.Background(), stmt, sketchId)
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
		values := []any{}
		for i, tag := range tagsToInsert {
			query += fmt.Sprintf("($1, $%d),", i+2)
			values = append(values, tag)
		}
		query = query[:len(query)-1] // Trim last comma
		values = append([]any{sketchId}, values...)
		fmt.Printf("QUERY: %s\n", query)
		fmt.Printf("VALUES: %+v\n", values)

		_, err = tx.Exec(context.Background(), query, values...)
		if err != nil {
			return err
		}
	}

	if len(tagsToDelete) > 0 {
		query := "DELETE FROM sketch_tags WHERE sketch_id = $1 AND tag_id IN ("
		values := []any{sketchId}
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
				websearch_to_tsquery('english', $%d)
				) AS rank
		`, args.ArgIndex)
		args.Args = append(args.Args, filter.Query)
	}

	baseFields := `
		 v.id as sketch_id, v.title as sketch_title, v.sketch_number as sketch_number, 
		v.slug as sketch_slug, v.thumbnail_name as thumbnail_name,
		v.upload_date as upload_date, v.rating as rating,
		c.id as creator_id, c.name as creator_name, c.slug as creator_slug, 
		c.profile_img as creator_img, sh.id as show_id, sh.name as show_name,
		sh.profile_img as show_img, sh.slug as show_slug, v.popularity_score as popularity,
		se.id as season_id, se.slug as season_slug, se.season_number as season_number, 
		e.id as episode_id, e.slug as episode_slug, e.episode_number as episode_number, e.air_date as episode_airdate,
		(select thumbnail_name from cast_members where %s and sketch_id = v.id order by position limit 1) as cast_thumbnail_name
		%s
	`

	castThumbnailClause := ""
	if len(filter.PersonIDs) == 0 && len(filter.CharacterIDs) == 0 {
		castThumbnailClause = "1=1"
	}
	if len(filter.PersonIDs) != 0 {
		personId := filter.PersonIDs[0]
		args.ArgIndex++
		args.Args = append(args.Args, personId)
		castThumbnailClause = fmt.Sprintf("person_id = $%d", args.ArgIndex)
	}
	if len(filter.CharacterIDs) != 0 {
		characterId := filter.CharacterIDs[0]
		args.ArgIndex++
		args.Args = append(args.Args, characterId)
		if castThumbnailClause == "" {
			castThumbnailClause = fmt.Sprintf("character_id = $%d", args.ArgIndex)
		} else {
			castThumbnailClause += fmt.Sprintf(" AND character_id = $%d", args.ArgIndex)
		}
	}

	fields := fmt.Sprintf(baseFields, castThumbnailClause, rankParam)

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
				COALESCE(c.alias, '') || ' ' || COALESCE(sh.name, '') || 
				' ' || COALESCE(sh.aliases,'') || ' ' ||
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
				), ' '),'')) @@ websearch_to_tsquery('english', $%d)
			`, args.ArgIndex)
		args.Args = append(args.Args, filter.Query)
	}

	// NOTE: Creators, shows and tags use OR operator
	if len(filter.CreatorIDs) > 0 {
		creatorPlaceholders := []string{}
		for _, creatorId := range filter.CreatorIDs {
			args.ArgIndex++
			creatorPlaceholders = append(creatorPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, creatorId)
		}

		clause += fmt.Sprintf(" AND c.id IN (%s)", strings.Join(creatorPlaceholders, ","))
	}

	if len(filter.CharacterIDs) > 0 {
		characterPlaceholders := []string{}
		for _, characterId := range filter.CharacterIDs {
			args.ArgIndex++
			characterPlaceholders = append(characterPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, characterId)
		}

		clause += fmt.Sprintf(" AND cm.character_id IN (%s)", strings.Join(characterPlaceholders, ","))
	}

	if len(filter.ShowIDs) > 0 {
		showPlaceholders := []string{}
		for _, showId := range filter.ShowIDs {
			args.ArgIndex++
			showPlaceholders = append(showPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, showId)
		}

		clause += fmt.Sprintf(" AND sh.id IN (%s)", strings.Join(showPlaceholders, ","))
	}

	if len(filter.TagIDs) > 0 {
		tagPlaceholders := []string{}
		for _, tagId := range filter.TagIDs {
			args.ArgIndex++
			tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, tagId)
		}

		clause += fmt.Sprintf(" AND vt.tag_id IN (%s)", strings.Join(tagPlaceholders, ","))
	}

	// NOTE: People filter use AND operation
	if len(filter.PersonIDs) > 0 {
		peoplePlaceholders := []string{}
		for _, personId := range filter.PersonIDs {
			args.ArgIndex++
			peoplePlaceholders = append(peoplePlaceholders, fmt.Sprintf("$%d", args.ArgIndex))
			args.Args = append(args.Args, personId)
		}

		clause += fmt.Sprintf(" AND cm.person_id IN (%s)", strings.Join(peoplePlaceholders, ","))
		clause += `
		GROUP BY v.id, v.title, v.slug,
		         v.thumbnail_name, v.upload_date, v.sketch_number,
		         c.id, c.name, c.page_url, c.slug, c.profile_img,
				sh.id, sh.name, sh.profile_img, sh.slug, 
				se.id, se.slug, se.season_number, 
				e.id, e.slug, e.episode_number, e.air_date`

		if filter.Query != "" {
			clause += ", rank"
		}

		if len(filter.PersonIDs) > 1 {
			args.ArgIndex++
			clause += fmt.Sprintf(" HAVING COUNT(DISTINCT cm.person_id) = $%d ", args.ArgIndex)
			args.Args = append(args.Args, len(filter.PersonIDs))
		}
	}

	return clause
}

func determineSort(filter *Filter, args *Arguements) string {
	sort := "upload_date ASC, popularity ASC"
	if val, ok := sortMap[filter.SortBy]; ok {
		sort = val
	}

	sort = fmt.Sprintf(" ORDER BY %s", sort)
	sort += fmt.Sprintf(" LIMIT $%d OFFSET $%d", args.ArgIndex+1, args.ArgIndex+2)
	args.ArgIndex += 2
	args.Args = append(args.Args, filter.Limit(), filter.Offset())

	return sort
}

func (m *SketchModel) Get(filter *Filter) ([]*SketchRef, Metadata, error) {
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
		SELECT count(*) OVER() as total_count, sketch_id, sketch_title, sketch_number,
		sketch_slug, thumbnail_name, upload_date, rating,
		creator_id, creator_name, creator_slug, creator_img, 
		show_id, show_name, show_img, show_slug, 
		season_id, season_slug, season_number, 
		episode_id, episode_slug, episode_number, episode_airdate,
		cast_thumbnail_name, popularity %s
		FROM ranked_sketches
		WHERE rn = 1
		%s
	`

	args := &Arguements{ArgIndex: 0}

	rank := ""
	if filter.Query != "" {
		rank = ", rank"
	}

	fields := determineFields(filter, args)
	conditionClause := determineConditions(filter, args)
	sortClause := determineSort(filter, args)

	// fmt.Println("SORT: ", sortClause)
	query = fmt.Sprintf(query, fields, conditionClause, rank, sortClause)
	fmt.Println(query)
	fmt.Printf("ARGS: %+v\n", args.Args)

	rows, err := m.DB.Query(context.Background(), query, args.Args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Metadata{}, ErrNoRecord
		} else {
			return nil, Metadata{}, err
		}
	}
	defer rows.Close()

	sketches := []*SketchRef{}
	var totalCount int

	for rows.Next() {
		v := &SketchRef{}
		c := &CreatorRef{}
		sh := &ShowRef{}
		se := &SeasonRef{}
		ep := &EpisodeRef{}
		destinations := []any{
			&totalCount, &v.ID, &v.Title, &v.Number, &v.Slug, &v.Thumbnail,
			&v.UploadDate, &v.Rating,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
			&se.ID, &se.Slug, &se.Number,
			&ep.ID, &ep.Slug, &ep.Number, &ep.AirDate,
			&v.CastThumbnail, nil,
		}
		var rank *float32
		if filter.Query != "" {
			destinations = append(destinations, &rank)
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, Metadata{}, err
		}

		se.Show = sh
		ep.Season = se

		if c.ID != nil {
			v.Creator = c
		}

		if ep.ID != nil {
			v.Episode = ep
		}
		sketches = append(sketches, v)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return sketches, calculateMetadata(totalCount, filter.Page, filter.PageSize), nil
}

func (m *SketchModel) GetById(id int) (*Sketch, error) {
	stmt := `
		SELECT v.id, v.title, v.sketch_number, v.sketch_url, v.description,
		v.slug, v.thumbnail_name, v.upload_date, v.youtube_id, v.popularity_score,
		v.episode_start, v.part_number, v.duration, v.rating, v.total_ratings,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		p.id, p.slug, p.first, p.last, p.profile_img,
		ch.id, ch.name, ch.slug, ch.img_name, ch.character_type,
		cm.id, cm.position, cm.character_name, cm.role, cm.profile_img, cm.thumbnail_name,
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name,
		se.id, se.slug, se.season_number,
		ser.id, ser.slug, ser.title, ser.thumbnail_name,
		rec.id, rec.slug, rec.title, rec.thumbnail_name
		FROM sketch AS v
		LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN cast_members as cm ON v.id = cm.sketch_id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN character as ch ON cm.character_id = ch.id
		LEFT JOIN series as ser ON v.series_id = ser.id
		LEFT JOIN recurring as rec ON v.recurring_id = rec.id
		WHERE v.id = $1
		ORDER BY cm.position asc, cm.id asc
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
	c := &CreatorRef{}
	sh := &ShowRef{}
	s := &SeasonRef{}
	e := &Episode{}
	rec := &RecurringRef{}
	se := &SeriesRef{}
	members := []*CastMember{}
	hasRows := false
	for rows.Next() {
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		hasRows = true
		err := rows.Scan(
			&v.ID, &v.Title, &v.Number, &v.URL, &v.Description, &v.Slug, &v.ThumbnailName,
			&v.UploadDate, &v.YoutubeID, &v.Popularity, &v.EpisodeStart, &v.SeriesPart,
			&v.Duration, &v.Rating, &v.TotalRatings,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&ch.ID, &ch.Name, &ch.Slug, &ch.Image, &ch.Type,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg, &cm.ThumbnailName,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&s.ID, &s.Slug, &s.Number,
			&se.ID, &se.Slug, &se.Title, &se.ThumbnailName,
			&rec.ID, &rec.Slug, &rec.Title, &rec.ThumbnailName,
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

	s.Show = sh
	e.Season = s
	if e.ID != nil {
		v.Episode = e
	}

	if c.ID != nil {
		v.Creator = c
	}

	v.Cast = members
	if se.ID != nil {
		v.Series = se
	}
	if rec.ID != nil {
		v.Recurring = rec
	}
	return v, nil
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

	args := &Arguements{}
	fields := determineFields(filter, args)
	conditionClause := determineConditions(filter, args)
	query = fmt.Sprintf(query, fields, conditionClause)
	// fmt.Println("COUNT QUERY", query)

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
		SELECT v.id, v.title, v.slug, v.thumbnail_name, 
			v.upload_date, v.rating,
			c.id, c.name, c.profile_img, c.slug,
			sh.id, sh.name, sh.profile_img, sh.slug,
			se.id, se.slug, se.season_number,
			e.id, e.slug, e.title, e.episode_number, e.air_date, e.thumbnail_name,	
			p.id, p.slug, p.first, p.last, p.profile_img,
			cm.id, cm.position, cm.thumbnail_name, cm.character_name,
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
		ORDER BY v.title, cm.position
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
		c := &CreatorRef{}
		sh := &ShowRef{}
		se := &SeasonRef{}
		ep := &Episode{}
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		hasRows = true

		err := rows.Scan(
			&v.ID, &v.Title, &v.Slug, &v.ThumbnailName,
			&v.UploadDate, &v.Rating,
			&c.ID, &c.Name, &c.ProfileImage, &c.Slug,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
			&se.ID, &se.Slug, &se.Number,
			&ep.ID, &ep.Slug, &ep.Title, &ep.Number, &ep.AirDate, &ep.Thumbnail,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&cm.ID, &cm.Position, &cm.ThumbnailName, &cm.CharacterName,
			&ch.ID, &ch.Slug, &ch.Name, &ch.Image,
		)
		if err != nil {
			return nil, err
		}

		se.Show = sh
		ep.Season = se
		if ep.ID != nil {
			v.Episode = ep
		}

		if c.ID != nil {
			v.Creator = c
		}
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

func (m *SketchModel) GetByUserLikes(userId int) ([]*SketchRef, error) {
	stmt := `
		SELECT v.id as sketch_id, v.title as sketch_title, v.sketch_number as sketch_number,
		v.slug as sketch_slug, v.thumbnail_name as thumbnail_name, 
		v.upload_date as upload_date, v.rating,
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

	sketches := []*SketchRef{}
	for rows.Next() {
		v := &SketchRef{}
		c := &CreatorRef{}
		sh := &ShowRef{}
		se := &SeasonRef{}
		ep := &EpisodeRef{}
		destinations := []any{
			&v.ID, &v.Title, &v.Number, &v.Slug, &v.Thumbnail, &v.UploadDate, &v.Rating,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug, &se.Number, &ep.Number,
		}
		err := rows.Scan(destinations...)
		if err != nil {
			return nil, err
		}
		v.Creator = c

		se.Show = sh
		ep.Season = se
		v.Episode = ep
		sketches = append(sketches, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sketches, nil
}

func (m *SketchModel) HasLike(sketchId, userId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM likes WHERE sketch_id = $1 AND user_id = $2)"
	err := m.DB.QueryRow(context.Background(), stmt, sketchId, userId).Scan(&exists)
	return exists, err
}

func (m *SketchModel) Insert(sketch *Sketch) (int, error) {
	var episodeId *int
	if sketch.Episode != nil {
		episodeId = sketch.Episode.ID
	}
	var seriesId *int
	if sketch.Series != nil {
		seriesId = sketch.Series.ID
	}

	var recurringId *int
	if sketch.Recurring != nil {
		recurringId = sketch.Recurring.ID
	}

	stmt := `
	INSERT INTO sketch (
		title, sketch_url, thumbnail_name, upload_date, slug, youtube_id, sketch_number,
		episode_id, episode_start, series_id, part_number, recurring_id, duration, description,
		popularity_score)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	RETURNING id;`
	result := m.DB.QueryRow(
		context.Background(), stmt, sketch.Title,
		sketch.URL, sketch.ThumbnailName, sketch.UploadDate,
		sketch.Slug, sketch.YoutubeID, sketch.Number, episodeId,
		sketch.EpisodeStart, seriesId, sketch.SeriesPart, recurringId,
		sketch.Duration, sketch.Description, sketch.Popularity,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	sketch.ID = &id
	return id, nil
}

func (m *SketchModel) InsertThumbnailName(sketchId int, name string) error {
	stmt := `UPDATE sketch SET thumbnail_name = $1 WHERE id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, name, sketchId)
	return err
}

func (m *SketchModel) InsertSketchCreatorRelation(sketchId, creatorId int) error {
	stmt := `INSERT INTO sketch_creator_rel (sketch_id, creator_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, sketchId, creatorId)
	return err
}

func (m *SketchModel) UpdateCreatorRelation(sketchId, creatorId int) error {
	stmt := `UPDATE sketch_creator_rel SET creator_id = $1 WHERE sketch_id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, creatorId, sketchId)
	return err
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

func (m *SketchModel) Update(sketch *Sketch) error {
	var episodeId *int
	if sketch.Episode != nil {
		episodeId = sketch.Episode.ID
	}
	var seriesId *int
	if sketch.Series != nil {
		seriesId = sketch.Series.ID
	}
	var recurringId *int
	if sketch.Recurring != nil {
		recurringId = sketch.Recurring.ID
	}

	stmt := `
	UPDATE sketch SET title = $1, sketch_url = $2, upload_date = $3, 
	slug = $4, thumbnail_name = $5, sketch_number = $6, episode_id = $7, episode_start = $8,
	series_id = $9, part_number = $10, recurring_id = $11, duration = $12, description = $13,
	popularity_score = $14, youtube_id = $15
	WHERE id = $16
	`
	_, err := m.DB.Exec(
		context.Background(), stmt,
		sketch.Title, sketch.URL, sketch.UploadDate, sketch.Slug, sketch.ThumbnailName,
		sketch.Number, episodeId, sketch.EpisodeStart, seriesId, sketch.SeriesPart,
		recurringId, sketch.Duration, sketch.Description, sketch.Popularity, sketch.YoutubeID,
		sketch.ID,
	)
	return err
}

func (m *SketchModel) SyncSketchCreators(sketchID int, creatorIDs []int) error {
	ctx := context.Background()
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if len(creatorIDs) == 0 {
		_, err = tx.Exec(ctx, `DELETE FROM sketch_creator_rel WHERE sketch_id = $1`, sketchID)
		if err != nil {
			return err
		}
		return tx.Commit(ctx)
	}

	// delete removed
	_, err = tx.Exec(ctx, `
		DELETE FROM sketch_creator_rel
		WHERE sketch_id = $1
		  AND creator_id <> ALL($2::int[])
	`, sketchID, creatorIDs)
	if err != nil {
		return err
	}

	// upsert present with order from list position
	_, err = tx.Exec(ctx, `
		WITH items AS (
		  SELECT x.creator_id, x.ord
		  FROM UNNEST($2::int[]) WITH ORDINALITY AS x(creator_id, ord)
		)
		INSERT INTO sketch_creator_rel (sketch_id, creator_id, position)
		SELECT $1, items.creator_id, items.ord
		FROM items
		ON CONFLICT (sketch_id, creator_id)
		DO UPDATE SET position = EXCLUDED.position
	`, sketchID, creatorIDs)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
