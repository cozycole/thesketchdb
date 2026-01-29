package models

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Moment struct {
	ID        *int
	Timestamp *int
	Sketch    *SketchRef
	Quotes    []*Quote
}

type Quote struct {
	ID         *int
	CastMember *CastMember
	Text       *string
	Type       *string
	Funny      *string
	Position   *int
	Moment     *Moment
	Tags       []*Tag
}

type MomentModelInterface interface {
	BatchUpdateQuotes(int, []*Quote) error
	BatchUpdateQuoteTags(int, []*Tag) error
	Delete(int) error
	Insert(int, *Moment) (int, error)
	GetById(int) (*Moment, error)
	GetBySketch(int) ([]*Moment, error)
	GetQuote(int) (*Quote, error)
	Update(*Moment) error
}

type MomentModel struct {
	DB *pgxpool.Pool
}

func (m *MomentModel) Delete(id int) error {
	stmt := `
		DELETE FROM moment
		WHERE id = $1
	`

	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}

func (m *MomentModel) Insert(sketchId int, moment *Moment) (int, error) {
	stmt := ` 
		INSERT INTO moment (sketch_id, timestamp) VALUES ($1,$2)
		RETURNING id;
	`

	result := m.DB.QueryRow(context.Background(), stmt, sketchId, moment.Timestamp)
	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *MomentModel) GetById(momentId int) (*Moment, error) {
	stmt := `
		SELECT m.id, m.timestamp, 
		q.id, q.text, q.type, q.funny, q.position,
		cm.id, cm.position, cm.character_name, cm.role, cm.profile_img, cm.thumbnail_name,
		s.id, s.title, s.sketch_number, s.slug, s.thumbnail_name, s.upload_date,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		se.id, se.slug, se.season_number,
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name, 
		t.id, t.name, t.slug, t.type,
		ca.id, ca.name, ca.slug
		FROM moment as m
		JOIN sketch as s ON m.sketch_id = s.id
		LEFT JOIN quote as q ON m.id = q.moment_id
		LEFT JOIN quote_tags_rel as qtr ON q.id = qtr.quote_id
		LEFT JOIN tags as t ON qtr.tag_id = t.id
		LEFT JOIN categories as ca ON t.category_id = ca.id
		LEFT JOIN cast_members as cm ON q.cast_id = cm.id
		LEFT JOIN sketch_creator_rel as vcr ON s.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON s.episode_id = e.id 
		LEFT JOIN season as se ON  e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		WHERE m.id = $1
		ORDER BY q.position
	`

	rows, err := m.DB.Query(context.Background(), stmt, momentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	s := &SketchRef{}
	c := &CreatorRef{}
	sh := &ShowRef{}
	se := &SeasonRef{}
	e := &EpisodeRef{}
	hasRows := false
	moment := &Moment{}
	quoteMap := map[int]*Quote{}
	quoteTagsMap := map[int]map[int]*Tag{}
	for rows.Next() {
		q := &Quote{}
		cm := &CastMember{}

		hasRows = true
		t := Tag{}
		ca := Category{}
		err := rows.Scan(
			&moment.ID, &moment.Timestamp,
			&q.ID, &q.Text, &q.Type, &q.Funny, &q.Position,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg,
			&cm.ThumbnailName,
			&s.ID, &s.Title, &s.Number, &s.Slug, &s.Thumbnail, &s.UploadDate,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&se.ID, &se.Slug, &se.Number,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&t.ID, &t.Name, &t.Slug, &t.Type,
			&ca.ID, &ca.Name, &ca.Slug,
		)

		if err != nil {
			return nil, err
		}

		se.Show = sh
		e.Season = se

		s.Episode = e
		s.Creator = c

		moment.Sketch = s

		if q.ID == nil {
			continue
		}

		q.CastMember = cm

		quoteMap[*q.ID] = q

		if t.ID != nil {
			t.Category = &ca
			if tm, ok := quoteTagsMap[*q.ID]; ok {
				tm[*t.ID] = &t
			} else {
				quoteTagsMap[*q.ID] = map[int]*Tag{*t.ID: &t}
			}
		}
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// create quotes slice for moment
	for _, q := range quoteMap {
		moment.Quotes = append(moment.Quotes, q)
	}

	sort.Slice(moment.Quotes, func(i, j int) bool {
		return *moment.Quotes[i].Position < *moment.Quotes[j].Position
	})

	// add tags to each quote
	for _, q := range moment.Quotes {
		if tagMap, ok := quoteTagsMap[*q.ID]; ok {
			for _, tag := range tagMap {
				q.Tags = append(q.Tags, tag)
			}
		}
	}

	return moment, nil
}

func (m *MomentModel) GetBySketch(sketchId int) ([]*Moment, error) {
	stmt := `
		SELECT m.id, m.timestamp, 
		q.id, q.text, q.type, q.funny, q.position,
		cm.id, cm.position, cm.character_name, cm.role, cm.profile_img, cm.thumbnail_name,
		p.id, p.slug, p.first, p.last, p.profile_img,
		ch.id, ch.name, ch.img_name,
		s.id, s.title, s.sketch_number,
		s.slug, s.thumbnail_name, s.upload_date,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		se.id, se.slug, se.season_number,
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name,
		t.id, t.name, t.slug, t.type,
		ca.id, ca.name, ca.slug
		FROM moment as m
		JOIN sketch as s ON m.sketch_id = s.id
		LEFT JOIN quote as q ON m.id = q.moment_id
		LEFT JOIN quote_tags_rel as qtr ON q.id = qtr.quote_id
		LEFT JOIN tags as t ON qtr.tag_id = t.id
		LEFT JOIN categories as ca ON t.category_id = ca.id
		LEFT JOIN cast_members as cm ON q.cast_id = cm.id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN character as ch ON cm.character_id = ch.id
		LEFT JOIN sketch_creator_rel as vcr ON s.id = vcr.sketch_id
		LEFT JOIN creator as c ON vcr.creator_id = c.id
		LEFT JOIN episode as e ON s.episode_id = e.id 
		LEFT JOIN season as se ON  e.season_id = se.id
		LEFT JOIN show as sh ON se.show_id = sh.id
		WHERE s.id = $1
		ORDER BY m.timestamp asc
	`

	rows, err := m.DB.Query(context.Background(), stmt, sketchId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	s := &SketchRef{}
	c := &CreatorRef{}
	sh := &ShowRef{}
	se := &SeasonRef{}
	e := &EpisodeRef{}
	hasRows := false
	momentMap := map[int]*Moment{}
	momentQuoteMap := map[int]map[int]*Quote{}
	quoteTagMap := map[int]map[int]*Tag{}
	for rows.Next() {
		m := &Moment{}
		q := &Quote{}
		cm := &CastMember{}
		t := &Tag{}
		ca := &Category{}
		ch := &Character{}
		p := &Person{}

		hasRows = true
		err := rows.Scan(
			&m.ID, &m.Timestamp,
			&q.ID, &q.Text, &q.Type, &q.Funny, &q.Position,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg,
			&cm.ThumbnailName,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&ch.ID, &ch.Name, &ch.Image,
			&s.ID, &s.Title, &s.Number, &s.Slug, &s.Thumbnail, &s.UploadDate,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&se.ID, &se.Slug, &se.Number,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail,
			&t.ID, &t.Slug, &t.Name, &t.Type,
			&ca.ID, &ca.Name, &ca.Slug,
		)

		if err != nil {
			return nil, err
		}

		se.Show = sh
		e.Season = se
		s.Episode = e

		s.Creator = c

		m.Sketch = s

		if m.ID == nil {
			continue
		}

		momentMap[*m.ID] = m

		if q.ID == nil {
			continue
		}

		cm.Actor = p
		cm.Character = ch
		q.CastMember = cm

		if t.ID != nil {
			t.Category = ca
			if currentQuoteMap, ok := quoteTagMap[*q.ID]; ok {
				currentQuoteMap[*t.ID] = t
			} else {
				quoteTagMap[*q.ID] = map[int]*Tag{*t.ID: t}
			}
		}

		// save map of moments to quote
		// this is for duplicates from tag joins
		if currentMomentMap, ok := momentQuoteMap[*m.ID]; ok {
			currentMomentMap[*q.ID] = q
		} else {
			momentQuoteMap[*m.ID] = map[int]*Quote{*q.ID: q}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	moments := []*Moment{}
	for _, m := range momentMap {
		for _, q := range momentQuoteMap[safeDeref(m.ID)] {
			for _, t := range quoteTagMap[safeDeref(q.ID)] {
				q.Tags = append(q.Tags, t)
			}
			m.Quotes = append(m.Quotes, q)
		}

		// Sort ascending by Position
		sort.Slice(m.Quotes, func(i, j int) bool {
			return *m.Quotes[i].Position < *m.Quotes[j].Position
		})
		moments = append(moments, m)
	}

	sort.Slice(moments, func(i, j int) bool {
		return *moments[i].Timestamp < *moments[j].Timestamp
	})

	return moments, nil
}

func (m *MomentModel) GetQuote(quoteId int) (*Quote, error) {
	stmt := `
		SELECT q.id, q.text, q.type, q.funny, q.position, q.moment_id,
		cm.id, cm.position, cm.character_name, cm.role, 
		cm.profile_img, cm.thumbnail_name,
		p.id, p.slug, p.first, p.last,
		t.id, t.name, t.slug, t.type,
		c.id, c.name, c.slug
		FROM quote as q
		LEFT JOIN cast_members as cm ON q.cast_id = cm.id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN quote_tags_rel as qtr ON q.id = qtr.quote_id
		LEFT JOIN tags as t ON qtr.tag_id = t.id
		LEFT JOIN categories as c ON t.category_id = c.id
		WHERE q.id = $1
	`

	rows, err := m.DB.Query(context.Background(), stmt, quoteId)
	if err != nil {
		return nil, err
	}
	mo := Moment{}
	q := Quote{}
	cm := CastMember{}
	p := Person{}
	tags := []*Tag{}
	hasRows := false
	for rows.Next() {
		hasRows = true

		t := Tag{}
		c := Category{}
		err := rows.Scan(
			&q.ID, &q.Text, &q.Type, &q.Funny, &q.Position, &mo.ID,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole,
			&cm.ProfileImg, &cm.ThumbnailName,
			&p.ID, &p.Slug, &p.First, &p.Last,
			&t.ID, &t.Name, &t.Slug, &t.Type,
			&c.ID, &c.Name, &c.Slug,
		)

		if err != nil {
			return nil, err
		}
		if t.ID != nil {
			tags = append(tags, &t)
		}
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	cm.Actor = &p
	q.CastMember = &cm
	q.Tags = tags
	q.Moment = &mo

	return &q, nil
}

func (m *MomentModel) Update(moment *Moment) error {
	stmt := `
		UPDATE moment SET timestamp = $1
		WHERE id = $2
	`
	_, err := m.DB.Exec(context.Background(), stmt, moment.Timestamp, moment.ID)
	return err
}

// UpdateQuotes handles the complete update of quotes for a moment
func (m *MomentModel) BatchUpdateQuotes(momentID int, quotes []*Quote) error {
	// Start transaction
	ctx := context.Background()
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	existingQuotes, err := getExistingQuotes(ctx, tx, momentID)
	if err != nil {
		return fmt.Errorf("failed to get existing quotes: %w", err)
	}

	existingIds := make(map[int]*Quote)
	for _, quote := range existingQuotes {
		if safeDeref(quote.ID) != 0 {
			existingIds[*quote.ID] = quote
		}
	}

	// Track which existing quotes are still needed
	updatedIds := make(map[int]bool)
	for _, q := range quotes {
		id := safeDeref(q.ID)
		updatedIds[id] = true

		if existingQuote, exists := existingIds[id]; exists {
			// Update existing quote
			err = updateQuote(ctx, tx, existingQuote.ID, q, momentID)
			if err != nil {
				return fmt.Errorf("failed to update quote with id %d: %w", q.ID, err)
			}
		} else {
			// Insert new quote
			_, err = insertQuote(ctx, tx, q, momentID)
			if err != nil {
				return fmt.Errorf("failed to insert quote with id %d: %w", q.ID, err)
			}
		}
	}

	// Delete quotes that are no longer in the form
	for id, existingQuote := range existingIds {
		if !updatedIds[id] {
			err = deleteQuote(ctx, tx, existingQuote.ID)
			if err != nil {
				return fmt.Errorf("failed to delete quote with id %d: %w", *existingQuote.ID, err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getExistingQuotes(ctx context.Context, tx pgx.Tx, momentID int) ([]*Quote, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, cast_id, text, type, funny, position 
		FROM quote 
		WHERE moment_id = $1 
		ORDER BY position`, momentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []*Quote
	for rows.Next() {
		quote := &Quote{}
		var castID int
		err := rows.Scan(&quote.ID, &castID, &quote.Text, &quote.Type, &quote.Funny, &quote.Position)
		if err != nil {
			return nil, err
		}

		// Create a basic CastMember with just the ID
		quote.CastMember = &CastMember{ID: &castID}
		quotes = append(quotes, quote)
	}

	return quotes, rows.Err()
}

func insertQuote(ctx context.Context, tx pgx.Tx, quote *Quote, momentID int) (int, error) {
	if quote.CastMember == nil || safeDeref(quote.CastMember.ID) == 0 {
		return 0, fmt.Errorf("quote insert error: undefined cast member id")
	}

	row := tx.QueryRow(ctx, `
		INSERT INTO quote (moment_id, cast_id, text, type, funny, position)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		momentID, *quote.CastMember.ID, quote.Text, quote.Type, quote.Funny, quote.Position)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func updateQuote(ctx context.Context, tx pgx.Tx, quoteID *int, quote *Quote, momentID int) error {
	if quote.CastMember == nil || safeDeref(quote.CastMember.ID) == 0 {
		return fmt.Errorf("quote insert error: undefined cast member id")
	}
	_, err := tx.Exec(ctx, `
		UPDATE quote 
		SET cast_id = $2, text = $3, type = $4, funny = $5, position = $6
		WHERE id = $1 AND moment_id = $7`,
		*quoteID, *quote.CastMember.ID, quote.Text, quote.Type, quote.Funny, quote.Position, momentID)
	return err
}

func deleteQuote(ctx context.Context, tx pgx.Tx, quoteID *int) error {
	_, err := tx.Exec(ctx, "DELETE FROM quote WHERE id = $1", *quoteID)
	return err
}

func (m *MomentModel) BatchUpdateQuoteTags(quoteID int, tags []*Tag) error {
	ctx := context.Background()
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get existing tag associations for this quote
	existingTagIDs, err := getExistingQuoteTags(ctx, tx, quoteID)
	if err != nil {
		return fmt.Errorf("failed to get existing quote tags: %w", err)
	}

	// Create maps for efficient lookup
	existingTagMap := make(map[int]bool)
	for _, tagID := range existingTagIDs {
		existingTagMap[tagID] = true
	}

	newTagMap := make(map[int]bool)
	var newTagIDs []int
	for _, tag := range tags {
		if tag.ID != nil {
			tagID := *tag.ID
			newTagMap[tagID] = true
			newTagIDs = append(newTagIDs, tagID)
		}
	}

	// Find tags to insert (in newTagIDs but not in existing)
	var tagsToInsert []int
	for _, tagID := range newTagIDs {
		if !existingTagMap[tagID] {
			tagsToInsert = append(tagsToInsert, tagID)
		}
	}

	// Find tags to delete (in existing but not in newTagIDs)
	var tagsToDelete []int
	for _, existingTagID := range existingTagIDs {
		if !newTagMap[existingTagID] {
			tagsToDelete = append(tagsToDelete, existingTagID)
		}
	}

	// Insert new tag associations
	if len(tagsToInsert) > 0 {
		err = insertQuoteTagAssociations(ctx, tx, quoteID, tagsToInsert)
		if err != nil {
			return fmt.Errorf("failed to insert quote tag associations: %w", err)
		}
	}

	// Delete removed tag associations
	if len(tagsToDelete) > 0 {
		err = deleteQuoteTagAssociations(ctx, tx, quoteID, tagsToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete quote tag associations: %w", err)
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getExistingQuoteTags(ctx context.Context, tx pgx.Tx, quoteID int) ([]int, error) {
	rows, err := tx.Query(ctx, `
		SELECT tag_id 
		FROM quote_tags_rel 
		WHERE quote_id = $1 
		ORDER BY tag_id`, quoteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tagIDs []int
	for rows.Next() {
		var tagID int
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}

	return tagIDs, rows.Err()
}

func insertQuoteTagAssociations(ctx context.Context, tx pgx.Tx, quoteID int, tagIDs []int) error {
	// Use batch insert for efficiency
	batch := &pgx.Batch{}

	for _, tagID := range tagIDs {
		batch.Queue("INSERT INTO quote_tags_rel (quote_id, tag_id) VALUES ($1, $2)", quoteID, tagID)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	// Execute all batched queries
	for i := range len(tagIDs) {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert quote_tag association for tag %d: %w", tagIDs[i], err)
		}
	}

	return nil
}

func deleteQuoteTagAssociations(ctx context.Context, tx pgx.Tx, quoteID int, tagIDs []int) error {
	// Use a single query with IN clause for efficiency
	query := `DELETE FROM quote_tags_rel WHERE quote_id = $1 AND tag_id = ANY($2)`

	_, err := tx.Exec(ctx, query, quoteID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to delete quote tag associations: %w", err)
	}

	return nil
}
