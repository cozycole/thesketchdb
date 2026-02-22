package models

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Quote struct {
	ID          *int
	Text        *string
	Type        *string
	Funny       *string
	StartTimeMs *int
	EndTimeMs   *int
	CastMembers []*CastMember
	Tags        []*Tag
}

type TranscriptLine struct {
	ID         *int    `json:"id"`
	LineNumber *int    `json:"lineNumber"`
	Text       *string `json:"text"`
	StartMs    *int    `json:"startMs"`
	EndMs      *int    `json:"endMs"`
}

type QuoteModelInterface interface {
	// BatchUpdateQuotes(int, []*Quote) error
	BatchUpdateQuoteTags(int, []*Tag) error
	GetBySketch(int) ([]*Quote, error)
	GetTranscriptBySketch(int) ([]*TranscriptLine, error)
	// Get(int) (*Quote, error)
}

type QuoteModel struct {
	DB *pgxpool.Pool
}

func (m *QuoteModel) GetBySketch(sketchId int) ([]*Quote, error) {
	stmt := `
		SELECT q.id, q.text, q.type, q.funny, q.start_time_ms, q.end_time_ms,
		cm.id, cm.position, cm.character_name, cm.role, cm.profile_img, cm.thumbnail_name,
		p.id, p.slug, p.first, p.last, p.profile_img,
		ch.id, ch.slug, ch.name, ch.img_name,
		t.id, t.slug, t.name, t.type,
		ca.id, ca.slug, ca.name
		FROM quote as q
		LEFT JOIN quote_tags_rel as qtr ON q.id = qtr.quote_id
		LEFT JOIN tags as t ON qtr.tag_id = t.id
		LEFT JOIN categories as ca ON t.category_id = ca.id
		LEFT JOIN quote_cast_rel as qc ON q.id = qc.quote_id
		LEFT JOIN cast_members as cm ON qc.cast_id = cm.id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN character as ch ON cm.character_id = ch.id
		WHERE q.sketch_id = $1
		ORDER BY q.start_time_ms asc, q.id
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

	hasRows := false
	quotes := map[int]*Quote{}
	quoteCastMap := map[int]map[int]*CastMember{}
	quoteTagMap := map[int]map[int]*Tag{}
	for rows.Next() {
		q := &Quote{}
		cm := &CastMember{}
		t := &Tag{}
		ca := &Category{}
		ch := &CharacterRef{}
		p := &PersonRef{}

		hasRows = true
		err := rows.Scan(
			&q.ID, &q.Text, &q.Type, &q.Funny, &q.StartTimeMs, &q.EndTimeMs,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg,
			&cm.ThumbnailName,
			&p.ID, &p.Slug, &p.First, &p.Last, &p.ProfileImg,
			&ch.ID, &ch.Slug, &ch.Name, &ch.Image,
			&t.ID, &t.Slug, &t.Name, &t.Type,
			&ca.ID, &ca.Slug, &ca.Name,
		)

		if q.ID == nil {
			continue
		}

		if storedQuote, ok := quotes[*q.ID]; ok {
			q = storedQuote
		} else {
			quotes[*q.ID] = q
		}

		if cm.ID != nil {
			if p.ID != nil {
				cm.Actor = p
			}
			if ch.ID != nil {
				cm.Character = ch
			}

			// we use map of maps to ensure no duplicates
			if cmMap, ok := quoteCastMap[*q.ID]; ok {
				cmMap[*cm.ID] = cm
			} else {
				quoteCastMap[*q.ID] = map[int]*CastMember{*cm.ID: cm}
			}
		}

		if t.ID != nil {
			if ca.ID != nil {
				t.Category = ca
			}

			if tm, ok := quoteTagMap[*q.ID]; ok {
				tm[*t.ID] = t
			} else {
				quoteTagMap[*q.ID] = map[int]*Tag{*t.ID: t}
			}
		}

		if err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	// iterate through quotes to add their respective
	// cast and tag lists
	quoteList := []*Quote{}
	for _, q := range quotes {
		q.CastMembers = []*CastMember{}
		cmMap, ok := quoteCastMap[*q.ID]
		if ok {
			for _, cm := range cmMap {
				q.CastMembers = append(q.CastMembers, cm)
			}
		}

		sort.Slice(q.CastMembers, func(i, j int) bool {
			return *q.CastMembers[i].Position < *q.CastMembers[j].Position
		})

		q.Tags = []*Tag{}
		tagMap, ok := quoteTagMap[*q.ID]
		if ok {
			for _, tag := range tagMap {
				q.Tags = append(q.Tags, tag)
			}
		}
		sort.Slice(q.Tags, func(i, j int) bool {
			return *q.Tags[i].ID < *q.Tags[j].ID
		})

		quoteList = append(quoteList, q)
	}

	sort.Slice(quoteList, func(i, j int) bool {
		return *quoteList[i].StartTimeMs < *quoteList[j].StartTimeMs
	})

	return quoteList, nil
}

func (m *QuoteModel) GetTranscriptBySketch(sketchId int) ([]*TranscriptLine, error) {
	stmt := `
		SELECT id, line_number, text, start_ms, end_ms
		FROM transcription_lines
		WHERE sketch_id = $1
		ORDER BY id
	`

	rows, err := m.DB.Query(context.Background(), stmt, sketchId)
	if err != nil {
		return nil, err
	}

	lines := []*TranscriptLine{}
	for rows.Next() {
		l := TranscriptLine{}
		err = rows.Scan(
			&l.ID, &l.LineNumber, &l.Text, &l.StartMs, &l.EndMs,
		)
		if err != nil {
			return nil, err
		}

		lines = append(lines, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// func (m *QuoteModel) GetQuote(quoteId int) (*Quote, error) {
// 	stmt := `
// 		SELECT q.id, q.text, q.type, q.funny, q.position, q.moment_id,
// 		cm.id, cm.position, cm.character_name, cm.role,
// 		cm.profile_img, cm.thumbnail_name,
// 		p.id, p.slug, p.first, p.last,
// 		t.id, t.name, t.slug, t.type,
// 		c.id, c.name, c.slug
// 		FROM quote as q
// 		LEFT JOIN cast_members as cm ON q.cast_id = cm.id
// 		LEFT JOIN person as p ON cm.person_id = p.id
// 		LEFT JOIN quote_tags_rel as qtr ON q.id = qtr.quote_id
// 		LEFT JOIN tags as t ON qtr.tag_id = t.id
// 		LEFT JOIN categories as c ON t.category_id = c.id
// 		WHERE q.id = $1
// 	`
//
// 	rows, err := m.DB.Query(context.Background(), stmt, quoteId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	mo := Moment{}
// 	q := Quote{}
// 	cm := CastMember{}
// 	p := PersonRef{}
// 	tags := []*Tag{}
// 	hasRows := false
// 	for rows.Next() {
// 		hasRows = true
//
// 		t := Tag{}
// 		c := Category{}
// 		err := rows.Scan(
// 			&q.ID, &q.Text, &q.Type, &q.Funny, &q.Position, &mo.ID,
// 			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole,
// 			&cm.ProfileImg, &cm.ThumbnailName,
// 			&p.ID, &p.Slug, &p.First, &p.Last,
// 			&t.ID, &t.Name, &t.Slug, &t.Type,
// 			&c.ID, &c.Name, &c.Slug,
// 		)
//
// 		if err != nil {
// 			return nil, err
// 		}
// 		if t.ID != nil {
// 			tags = append(tags, &t)
// 		}
// 	}
//
// 	if !hasRows {
// 		return nil, ErrNoRecord
// 	}
//
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	cm.Actor = &p
// 	q.CastMember = &cm
// 	q.Tags = tags
// 	q.Moment = &mo
//
// 	return &q, nil
// }

// UpdateQuotes handles the complete update of quotes for a moment
// func (m *QuoteModel) BatchUpdateQuotes(sketchID int, quotes []*Quote) error {
// 	// Start transaction
// 	ctx := context.Background()
// 	tx, err := m.DB.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	defer tx.Rollback(ctx)
//
// 	existingQuotes, err := getExistingQuotes(ctx, tx, sketchID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get existing quotes: %w", err)
// 	}
//
// 	existingIds := make(map[int]*Quote)
// 	for _, quote := range existingQuotes {
// 		if safeDeref(quote.ID) != 0 {
// 			existingIds[*quote.ID] = quote
// 		}
// 	}
//
// 	// Track which existing quotes are still needed
// 	updatedIds := make(map[int]bool)
// 	for _, q := range quotes {
// 		id := safeDeref(q.ID)
// 		updatedIds[id] = true
//
// 		if existingQuote, exists := existingIds[id]; exists {
// 			// Update existing quote
// 			err = updateQuote(ctx, tx, existingQuote.ID, q, sketchID)
// 			if err != nil {
// 				return fmt.Errorf("failed to update quote with id %d: %w", q.ID, err)
// 			}
// 		} else {
// 			// Insert new quote
// 			_, err = insertQuote(ctx, tx, q, sketchID)
// 			if err != nil {
// 				return fmt.Errorf("failed to insert quote with id %d: %w", q.ID, err)
// 			}
// 		}
// 	}
//
// 	// Delete quotes that are no longer in the form
// 	for id, existingQuote := range existingIds {
// 		if !updatedIds[id] {
// 			err = deleteQuote(ctx, tx, existingQuote.ID)
// 			if err != nil {
// 				return fmt.Errorf("failed to delete quote with id %d: %w", *existingQuote.ID, err)
// 			}
// 		}
// 	}
//
// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
//
// 	return nil
// }
//
// func getExistingQuotes(ctx context.Context, tx pgx.Tx, momentID int) ([]*Quote, error) {
// 	rows, err := tx.Query(ctx, `
// 		SELECT id, cast_id, text, type, funny, position
// 		FROM quote
// 		WHERE moment_id = $1
// 		ORDER BY position`, momentID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var quotes []*Quote
// 	for rows.Next() {
// 		quote := &Quote{}
// 		var castID int
// 		err := rows.Scan(&quote.ID, &castID, &quote.Text, &quote.Type, &quote.Funny, &quote.Position)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Create a basic CastMember with just the ID
// 		quotes = append(quotes, quote)
// 	}
//
// 	return quotes, rows.Err()
// }

func insertQuote(ctx context.Context, tx pgx.Tx, quote *Quote) (int, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO quote (start_time_ms, end_time_ms, text, type, funny)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		quote.StartTimeMs, quote.EndTimeMs, quote.Text, quote.Type, quote.Funny)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func updateQuote(ctx context.Context, tx pgx.Tx, quote *Quote) error {
	if safeDeref(quote.ID) == 0 {
		return fmt.Errorf("quote insert error: undefined quote id")
	}

	_, err := tx.Exec(ctx, `
		UPDATE quote 
		SET start_time_ms = $2, end_time_ms = $3, text = $4, type = $5, funny = $6,
		WHERE id = $1`,
		quote.ID, quote.StartTimeMs, quote.EndTimeMs, quote.Text, quote.Type, quote.Funny)
	return err
}

func deleteQuote(ctx context.Context, tx pgx.Tx, quoteID *int) error {
	_, err := tx.Exec(ctx, "DELETE FROM quote WHERE id = $1", *quoteID)
	return err
}

func (m *QuoteModel) BatchUpdateQuoteTags(quoteID int, tags []*Tag) error {
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
