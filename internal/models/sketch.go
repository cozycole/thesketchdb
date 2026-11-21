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
	ID                *int          `json:"id"`
	Slug              *string       `json:"slug"`
	Title             *string       `json:"title"`
	URL               *string       `json:"url"`
	Duration          *int          `json:"duration"`
	VideoStartSeconds *int          `json:"videoStart"`
	Description       *string       `json:"description"`
	YoutubeID         *string       `json:"youtubeId"`
	ThumbnailName     *string       `json:"thumbnailName"`
	Popularity        *float32      `json:"popularity"`
	UploadDate        *Date         `json:"uploadDate"`
	Creator           *CreatorRef   `json:"creator"`
	Show              *ShowRef      `json:"show"`
	Episode           *Episode      `json:"episode"`
	EpisodeStart      *int          `json:"episodeStart"`
	Grouping          *Grouping     `json:"grouping,omitempty"`
	Cast              []*CastMember `json:"cast"`
	CastThumbnail     *string       `json:"castThumbnail"`
	Tags              []*Tag        `json:"tags"`
	Number            *int          `json:"episodeSketchOrder"`
	Series            *SeriesRef    `json:"series"`
	SeriesPart        *int          `json:"seriesPart"`
	Recurring         *RecurringRef `json:"recurring"`
	Rating            *float32      `json:"rating"`
	TotalRatings      *int          `json:"totalRatings"`
	Liked             *bool         `json:"liked,omitempty"`
}

type SketchRef struct {
	ID            *int        `json:"id"`
	Slug          *string     `json:"slug"`
	Title         *string     `json:"title"`
	Thumbnail     *string     `json:"thumbnail"`
	CastThumbnail *string     `json:"castThumbnail"`
	UploadDate    *Date       `json:"uploadDate"`
	Creator       *CreatorRef `json:"creator"`
	Show          *ShowRef    `json:"show"`
	Episode       *EpisodeRef `json:"episode"`
	Number        *int        `json:"episodeSketchOrder"`
	Rating        *float32    `json:"rating"`
}

type SketchVideo struct {
	ID           *int           `json:"id"`
	HotS3Key     *string        `json:"hotS3Key"`
	ColdS3Key    *string        `json:"coldS3Key"`
	ArchivedAt   *time.Time     `json:"archivedAt"`
	PipelineJobs []*PipelineJob `json:"jobs"`
}

type SketchModelInterface interface {
	BatchUpdateTags(sketchId int, tags []*Tag) error
	Delete(id int) error
	DeleteScreenshots(id int) error
	Exists(id int) (bool, error)
	Get(filter *Filter) ([]*SketchRef, Metadata, error)
	GetById(id int) (*Sketch, error)
	GetByUserLikes(id int) ([]*SketchRef, error)
	GetCount(filter *Filter) (int, error)
	GetFeatured() ([]*Sketch, error)
	GetVideos(int) ([]*SketchVideo, error)
	HasLike(sketchId, userId int) (bool, error)
	Insert(sketch *Sketch) (int, error)
	InsertSketchCreatorRelation(sketchId, creatorId int) error
	InsertSketchVideo(sketchId int, video *SketchVideo) error
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

func (m *SketchModel) DeleteScreenshots(id int) error {
	stmt := `
		DELETE from cast_auto_screenshots
		WHERE sketch_id = $1
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

func (m *SketchModel) BatchUpdateTags(sketchId int, tags []*Tag) error {
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

	// fmt.Printf("EXISTING TAGS: %+v\n", existingTags)

	// New tags
	newTags := make(map[int]bool)
	for _, tag := range tags {
		newTags[*tag.ID] = true
	}

	tagsToInsert := []int{}
	for tag_id := range newTags {
		if !existingTags[tag_id] {
			tagsToInsert = append(tagsToInsert, tag_id)
		}
	}
	// fmt.Printf("TAGS TO INSERT: %+v\n", tagsToInsert)

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
		// fmt.Printf("QUERY: %s\n", query)
		// fmt.Printf("VALUES: %+v\n", values)

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

var sketchListSorts = map[string]string{
	"popular": "v.popularity_score DESC, v.upload_date DESC, v.id DESC",
	"recent":  "v.id DESC",
	"newest":  "v.upload_date DESC, v.title ASC, v.id DESC",
	"oldest":  "v.upload_date ASC, v.title ASC, v.id DESC",
	"az":      "v.title ASC, v.id DESC",
	"za":      "v.title DESC, v.id DESC",
}

func sketchListOrder(sortBy string) string {
	if order, ok := sketchListSorts[sortBy]; ok {
		return order
	}
	return "v.upload_date ASC, v.popularity_score ASC, v.id DESC"
}

type sketchListBuilder struct {
	args           []any
	conditions     []string
	creatorWhere   string
	thumbnailWhere string
}

func (b *sketchListBuilder) bind(value any) string {
	b.args = append(b.args, value)
	return fmt.Sprintf("$%d", len(b.args))
}

func uniqueSketchFilterIDs(ids []int) []int32 {
	result := make([]int32, 0, len(ids))
	seen := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		if _, exists := seen[id]; !exists {
			seen[id] = struct{}{}
			result = append(result, int32(id))
		}
	}
	return result
}

// Each filter is independent. Values within each filter are ORed together,
// except PersonIDs, for which every requested person must appear in the sketch.
func (b *sketchListBuilder) addFilters(f Filter) {
	if q := strings.TrimSpace(f.Query); q != "" {
		param := b.bind(q)
		b.conditions = append(b.conditions, fmt.Sprintf(`
            EXISTS (
                SELECT 1 FROM sketch_search_documents sd
                WHERE sd.sketch_id = v.id
                  AND sd.document @@ websearch_to_tsquery('english', %s)
            )`, param))
	}

	if len(f.SketchIDs) > 0 {
		param := b.bind(uniqueSketchFilterIDs(f.SketchIDs))
		b.conditions = append(b.conditions,
			fmt.Sprintf("v.id = ANY(%s::int[])", param))
	}

	if len(f.CreatorIDs) > 0 {
		param := b.bind(uniqueSketchFilterIDs(f.CreatorIDs))
		b.conditions = append(b.conditions, fmt.Sprintf(`
            EXISTS (
                SELECT 1 FROM sketch_creator_rel scr
                WHERE scr.sketch_id = v.id
                  AND scr.creator_id = ANY(%s::int[])
            )`, param))
		// The displayed creator should be one of the requested creators.
		b.creatorWhere = fmt.Sprintf(" AND c.id = ANY(%s::int[])", param)
	}

	if len(f.ShowIDs) > 0 {
		param := b.bind(uniqueSketchFilterIDs(f.ShowIDs))
		b.conditions = append(b.conditions, fmt.Sprintf(`
            EXISTS (
                SELECT 1
                FROM episode e
                JOIN season se ON se.id = e.season_id
                WHERE e.id = v.episode_id
                  AND se.show_id = ANY(%s::int[])
            )`, param))
	}

	if len(f.TagIDs) > 0 {
		param := b.bind(uniqueSketchFilterIDs(f.TagIDs))
		b.conditions = append(b.conditions, fmt.Sprintf(`
            EXISTS (
                SELECT 1 FROM sketch_tags st
                WHERE st.sketch_id = v.id
                  AND st.tag_id = ANY(%s::int[])
            )`, param))
	}

	// The original query applies character and person conditions to the SAME
	// cast member rows. Keep that behavior when both filters are supplied.
	if len(f.PersonIDs) > 0 {
		ids := uniqueSketchFilterIDs(f.PersonIDs)
		param := b.bind(ids)
		castWhere := fmt.Sprintf("cm.person_id = ANY(%s::int[])", param)
		if len(f.CharacterIDs) > 0 {
			characters := b.bind(uniqueSketchFilterIDs(f.CharacterIDs))
			castWhere += fmt.Sprintf(" AND cm.character_id = ANY(%s::int[])", characters)
		}
		count := b.bind(len(ids))
		b.conditions = append(b.conditions, fmt.Sprintf(`
            v.id IN (
                SELECT cm.sketch_id
                FROM cast_members cm
                WHERE %s
                GROUP BY cm.sketch_id
                HAVING COUNT(DISTINCT cm.person_id) = %s
            )`, castWhere, count))
	} else if len(f.CharacterIDs) > 0 {
		param := b.bind(uniqueSketchFilterIDs(f.CharacterIDs))
		b.conditions = append(b.conditions, fmt.Sprintf(`
            EXISTS (
                SELECT 1 FROM cast_members cm
                WHERE cm.sketch_id = v.id
                  AND cm.character_id = ANY(%s::int[])
            )`, param))
	}

	// Filter.Type was not used in the old query. Its database meaning is not
	// defined in the supplied schema, so it is intentionally left unhandled.
	// Add a new condition here once the intended values/column are specified.
}

func (b *sketchListBuilder) addThumbnailFilter(f Filter) {
	// Match the existing behavior: use the FIRST selected person and character
	// for the optional cast-specific thumbnail, rather than all selected IDs.
	if len(f.PersonIDs) > 0 {
		b.thumbnailWhere += fmt.Sprintf(" AND cm.person_id = %s", b.bind(f.PersonIDs[0]))
	}
	if len(f.CharacterIDs) > 0 {
		b.thumbnailWhere += fmt.Sprintf(" AND cm.character_id = %s", b.bind(f.CharacterIDs[0]))
	}
}

const sketchListSQL = `
	SELECT
		page.total_count,
		v.id, v.title, v.sketch_number, v.slug, v.thumbnail_name,
		v.upload_date, v.rating,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.profile_img, sh.slug,
		se.id, se.slug, se.season_number,
		e.id, e.slug, e.episode_number, e.air_date,
		shg.id, shg.slug, shg.name, shg.profile_img,
		portrait.thumbnail_name
	FROM (
		-- There is exactly ONE row per sketch at this stage.
		SELECT v.id, COUNT(*) OVER () AS total_count
		FROM sketch v
		WHERE %s
		ORDER BY %s
		LIMIT %s OFFSET %s
	) AS page
	JOIN sketch v ON v.id = page.id
	LEFT JOIN LATERAL (
		-- A sketch can have multiple creators; select one deterministically.
		SELECT c.id, c.name, c.slug, c.profile_img
		FROM sketch_creator_rel scr
		JOIN creator c ON c.id = scr.creator_id
		WHERE scr.sketch_id = v.id %s
		ORDER BY c.id
		LIMIT 1
	) c ON TRUE
	LEFT JOIN episode e ON e.id = v.episode_id
	LEFT JOIN season se ON se.id = e.season_id
	LEFT JOIN show sh ON sh.id = se.show_id
	LEFT JOIN sketch_grouping g ON g.id = v.grouping_id
	LEFT JOIN show shg ON shg.id = g.show_id
	LEFT JOIN LATERAL (
		SELECT cm.thumbnail_name
		FROM cast_members cm
		WHERE cm.sketch_id = v.id %s
		ORDER BY cm.position NULLS LAST, cm.id
		LIMIT 1
	) portrait ON TRUE
	ORDER BY %s
`

func (m *SketchModel) Get(filter *Filter) ([]*SketchRef, Metadata, error) {
	f := Filter{}
	if filter != nil {
		f = *filter
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	pageSize := f.Limit()
	offset := (page - 1) * pageSize

	b := &sketchListBuilder{}
	b.addFilters(f)
	where := "TRUE"
	if len(b.conditions) > 0 {
		where = strings.Join(b.conditions, " AND ")
	}

	// Retain just filter arguments for the out-of-range-page count fallback.
	filterArgs := append([]any(nil), b.args...)
	limitParam := b.bind(pageSize)
	offsetParam := b.bind(offset)
	b.addThumbnailFilter(f)
	order := sketchListOrder(f.SortBy)

	query := fmt.Sprintf(sketchListSQL,
		where, order, limitParam, offsetParam,
		b.creatorWhere, b.thumbnailWhere, order,
	)

	ctx := context.Background() // Prefer passing r.Context() in a future signature.
	rows, err := m.DB.Query(ctx, query, b.args...)
	if err != nil {
		return nil, Metadata{}, fmt.Errorf("list sketches: %w", err)
	}
	defer rows.Close()

	sketches := make([]*SketchRef, 0, pageSize)
	var totalCount int

	for rows.Next() {
		v := &SketchRef{}
		c := &CreatorRef{}
		sh := &ShowRef{}
		shg := &ShowRef{}
		se := &SeasonRef{}
		ep := &EpisodeRef{}

		err := rows.Scan(
			&totalCount,
			&v.ID, &v.Title, &v.Number, &v.Slug, &v.Thumbnail,
			&v.UploadDate, &v.Rating,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
			&se.ID, &se.Slug, &se.Number,
			&ep.ID, &ep.Slug, &ep.Number, &ep.AirDate,
			&shg.ID, &shg.Slug, &shg.Name, &shg.ProfileImg,
			&v.CastThumbnail,
		)
		if err != nil {
			return nil, Metadata{}, fmt.Errorf("scan listed sketch: %w", err)
		}

		se.Show = sh
		ep.Season = se
		if c.ID != nil {
			v.Creator = c
		}
		if shg.ID != nil {
			// Preserve the old list's grouping-show behavior.
			v.Show = shg
		}
		if ep.ID != nil {
			v.Episode = ep
		}
		sketches = append(sketches, v)
	}
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, fmt.Errorf("iterate listed sketches: %w", err)
	}
	rows.Close() // Release the connection before running the count fallback.

	if len(sketches) == 0 {
		// COUNT(*) OVER() returns no rows when OFFSET is past the last page.
		// Recover an accurate count in that case (including an empty catalog).
		countSQL := "SELECT COUNT(*) FROM sketch v WHERE " + where
		if err := m.DB.QueryRow(ctx, countSQL, filterArgs...).Scan(&totalCount); err != nil {
			return nil, Metadata{}, fmt.Errorf("count sketches: %w", err)
		}
	}

	return sketches, calculateMetadata(totalCount, page, pageSize), nil
}

func (m *SketchModel) GetById(id int) (*Sketch, error) {
	stmt := `
		SELECT v.id, v.title, v.sketch_number, v.sketch_url, v.description,
		v.slug, v.thumbnail_name, v.upload_date, v.youtube_id, v.popularity_score,
		v.episode_start, v.part_number, v.duration, v.rating, v.total_ratings, v.video_start_seconds,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		g.id, g.slug, g.title,
		shg.id, shg.name, shg.slug, shg.profile_img,
		p.id, p.slug, p.first, p.last, p.profile_img,
		ch.id, ch.name, ch.slug, ch.img_name, ch.character_type,
		cm.id, cm.position, cm.character_name, cm.role, cm.profile_img, cm.thumbnail_name,
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name, e.youtube_id,
		se.id, se.slug, se.season_number,
		ser.id, ser.slug, ser.title, ser.thumbnail_name,
		rec.id, rec.slug, rec.title, rec.thumbnail_name
		FROM sketch AS v
		LEFT JOIN sketch_creator_rel as vcr ON v.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON v.episode_id = e.id
		LEFT JOIN season as se ON e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		LEFT JOIN sketch_grouping as g ON v.grouping_id = g.id
		LEFT JOIN show as shg ON g.show_id = shg.id
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
	g := &Grouping{}
	// show for grouping
	shg := &ShowRef{}
	s := &SeasonRef{}
	e := &Episode{}
	rec := &RecurringRef{}
	se := &SeriesRef{}
	members := []*CastMember{}
	hasRows := false
	for rows.Next() {
		p := &PersonRef{}
		ch := &CharacterRef{}
		cm := &CastMember{}
		hasRows = true
		err := rows.Scan(
			&v.ID, &v.Title, &v.Number, &v.URL, &v.Description, &v.Slug, &v.ThumbnailName,
			&v.UploadDate, &v.YoutubeID, &v.Popularity, &v.EpisodeStart, &v.SeriesPart,
			&v.Duration, &v.Rating, &v.TotalRatings, &v.VideoStartSeconds,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&g.ID, &g.Slug, &g.Title,
			&shg.ID, &shg.Name, &shg.Slug, &shg.ProfileImg,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&ch.ID, &ch.Name, &ch.Slug, &ch.Image, &ch.Type,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg, &cm.ThumbnailName,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail, &e.YoutubeID,
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

	if sh.ID != nil {
		s.Show = sh
		v.Show = sh
	} else {
		s.Show = shg
		v.Show = shg
	}

	if g.ID != nil {
		v.Grouping = g
	}

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
	f := Filter{}
	if filter != nil {
		f = *filter
	}

	b := &sketchListBuilder{}
	b.addFilters(f)

	where := "TRUE"
	if len(b.conditions) > 0 {
		where = strings.Join(b.conditions, " AND ")
	}

	query := "SELECT COUNT(*) FROM sketch v WHERE " + where

	var count int
	err := m.DB.QueryRow(
		context.Background(),
		query,
		b.args...,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("count sketches: %w", err)
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
		p := &PersonRef{}
		ch := &CharacterRef{}
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
		c.profile_img as creator_img, 
		sh.id as show_id, sh.name as show_name, sh.profile_img as show_img, sh.slug as show_slug, 
		se.id, se.season_number as season_number, 
		e.id, e.episode_number as episode_number
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
			&sh.ID, &sh.Name, &sh.ProfileImg, &sh.Slug,
			&se.ID, &se.Number,
			&ep.ID, &ep.Number,
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

// NOTE: This function doesn't currently account for multiple videos assigned
// to a single sketch! Implement later...
func (m *SketchModel) GetVideos(sketchId int) ([]*SketchVideo, error) {
	stmt := `
		SELECT v.id, v.hot_s3_key, v.cold_s3_key, v.archived_at,
		p.id, p.status, p.error
		FROM sketch_video as v
		LEFT JOIN pipeline_jobs as p ON v.id = p.video_id
		WHERE v.sketch_id = $1
		ORDER BY v.id asc
		LIMIT 1
	`

	row := m.DB.QueryRow(context.Background(), stmt, sketchId)

	v := SketchVideo{}
	p := PipelineJob{}
	err := row.Scan(&v.ID, &v.HotS3Key, &v.ColdS3Key, &v.ArchivedAt,
		&p.ID, &p.Status, &p.Error)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	videos := []*SketchVideo{}
	if v.ID != nil {
		v.PipelineJobs = []*PipelineJob{&p}
		videos = append(videos, &v)
	}

	return videos, nil
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

	var groupingId *int
	if sketch.Grouping != nil {
		groupingId = sketch.Grouping.ID
	}

	stmt := `
	INSERT INTO sketch (
		title, sketch_url, thumbnail_name, upload_date, slug, youtube_id, sketch_number,
		episode_id, episode_start, series_id, part_number, recurring_id, duration, description,
		popularity_score, grouping_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	RETURNING id;`
	result := m.DB.QueryRow(
		context.Background(), stmt, sketch.Title,
		sketch.URL, sketch.ThumbnailName, sketch.UploadDate,
		sketch.Slug, sketch.YoutubeID, sketch.Number, episodeId,
		sketch.EpisodeStart, seriesId, sketch.SeriesPart, recurringId,
		sketch.Duration, sketch.Description, sketch.Popularity,
		groupingId,
	)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	sketch.ID = &id
	return id, nil
}

func (m *SketchModel) InsertSketchCreatorRelation(sketchId, creatorId int) error {
	stmt := `INSERT INTO sketch_creator_rel (sketch_id, creator_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, sketchId, creatorId)
	return err
}

func (m *SketchModel) InsertSketchVideo(sketchId int, video *SketchVideo) error {
	stmt := `
		INSERT INTO sketch_video (sketch_id, hot_s3_key, cold_s3_key, archived_at)
		VALUES ($1,$2,$3,$4)
		RETURNING id;
	`

	row := m.DB.QueryRow(context.Background(), stmt, sketchId,
		video.HotS3Key, video.ColdS3Key, video.ArchivedAt)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	video.ID = &id

	return nil
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

	var groupingId *int
	if sketch.Grouping != nil {
		groupingId = sketch.Grouping.ID
	}

	stmt := `
	UPDATE sketch SET title = $1, sketch_url = $2, upload_date = $3, 
	slug = $4, thumbnail_name = $5, sketch_number = $6, episode_id = $7, episode_start = $8,
	series_id = $9, part_number = $10, recurring_id = $11, duration = $12, description = $13,
	popularity_score = $14, youtube_id = $15, grouping_id = $16, video_start_seconds = $17
	WHERE id = $18
	`
	_, err := m.DB.Exec(
		context.Background(), stmt,
		sketch.Title, sketch.URL, sketch.UploadDate,
		sketch.Slug, sketch.ThumbnailName, sketch.Number, episodeId, sketch.EpisodeStart,
		seriesId, sketch.SeriesPart, recurringId, sketch.Duration, sketch.Description,
		sketch.Popularity, sketch.YoutubeID, groupingId, sketch.VideoStartSeconds,
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
