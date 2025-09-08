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
	Sketch    *Sketch
	Quotes    []*Quote
}

type Quote struct {
	ID         *int
	CastMember *CastMember
	Text       *string
	Type       *string
	Funny      *string
	Position   *int
	Tags       []*QuoteTag
}

type QuoteTag struct {
	ID   *int
	Text *string
}

type MomentModelInterface interface {
	BatchUpdateQuotes(int, []*Quote) error
	Delete(int) error
	Insert(int, *Moment) (int, error)
	GetById(int) (*Moment, error)
	GetBySketch(int) ([]*Moment, error)
	// GetQuotes(int) ([]*Quote, error)
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
		s.id, s.title, s.sketch_number, s.sketch_url, 
		s.slug, s.thumbnail_name, s.upload_date, s.youtube_id, s.popularity_score,
		s.part_number, s.duration,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		se.id, se.slug, se.season_number,
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name, e.youtube_id
		FROM moment as m
		JOIN sketch as s ON m.sketch_id = s.id
		LEFT JOIN quote as q ON m.id = q.moment_id
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

	s := &Sketch{}
	c := &Creator{}
	sh := &Show{}
	se := &Season{}
	e := &Episode{}
	hasRows := false
	moment := &Moment{}
	for rows.Next() {
		q := &Quote{}
		cm := &CastMember{}

		hasRows = true
		err := rows.Scan(
			&moment.ID, &moment.Timestamp,
			&q.ID, &q.Text, &q.Type, &q.Funny, &q.Position,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg,
			&cm.ThumbnailName,
			&s.ID, &s.Title, &s.Number, &s.URL, &s.Slug, &s.ThumbnailName,
			&s.UploadDate, &s.YoutubeID, &s.Popularity, &s.SeriesPart, &s.Duration,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&se.ID, &se.Slug, &se.Number,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail, &e.YoutubeID,
		)

		if err != nil {
			return nil, err
		}

		e.Show = sh
		e.Season = se

		s.Episode = e
		s.Season = se
		s.Show = sh
		s.Creator = c

		moment.Sketch = s

		if q.ID == nil {
			continue
		}

		q.CastMember = cm

		moment.Quotes = append(moment.Quotes, q)
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	return moment, nil
}

func (m *MomentModel) GetBySketch(sketchId int) ([]*Moment, error) {
	stmt := `
		SELECT m.id, m.timestamp, 
		q.id, q.text, q.type, q.funny, q.position,
		cm.id, cm.position, cm.character_name, cm.role, cm.profile_img, cm.thumbnail_name,
		s.id, s.title, s.sketch_number, s.sketch_url, 
		s.slug, s.thumbnail_name, s.upload_date, s.youtube_id, s.popularity_score,
		s.part_number, s.duration,
		c.id, c.name, c.slug, c.profile_img,
		sh.id, sh.name, sh.slug, sh.profile_img,
		se.id, se.slug, se.season_number,
		e.id, e.slug, e.episode_number, e.title, e.air_date, e.thumbnail_name, e.youtube_id
		FROM moment as m
		JOIN sketch as s ON m.sketch_id = s.id
		LEFT JOIN quote as q ON m.id = q.moment_id
		LEFT JOIN cast_members as cm ON q.cast_id = cm.id
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

	s := &Sketch{}
	c := &Creator{}
	sh := &Show{}
	se := &Season{}
	e := &Episode{}
	hasRows := false
	momentMap := map[int]*Moment{}
	momentQuoteMap := map[int]map[int]*Quote{}
	for rows.Next() {
		m := &Moment{}
		q := &Quote{}
		cm := &CastMember{}

		hasRows = true
		err := rows.Scan(
			&m.ID, &m.Timestamp,
			&q.ID, &q.Text, &q.Type, &q.Funny, &q.Position,
			&cm.ID, &cm.Position, &cm.CharacterName, &cm.CastRole, &cm.ProfileImg,
			&cm.ThumbnailName,
			&s.ID, &s.Title, &s.Number, &s.URL, &s.Slug, &s.ThumbnailName,
			&s.UploadDate, &s.YoutubeID, &s.Popularity, &s.SeriesPart, &s.Duration,
			&c.ID, &c.Name, &c.Slug, &c.ProfileImage,
			&sh.ID, &sh.Name, &sh.Slug, &sh.ProfileImg,
			&se.ID, &se.Slug, &se.Number,
			&e.ID, &e.Slug, &e.Number, &e.Title, &e.AirDate, &e.Thumbnail, &e.YoutubeID,
		)

		if err != nil {
			return nil, err
		}

		e.Show = sh
		e.Season = se

		s.Episode = e
		s.Season = se
		s.Show = sh
		s.Creator = c

		m.Sketch = s

		momentMap[safeDeref(m.ID)] = m

		if q.ID == nil {
			continue
		}

		q.CastMember = cm

		// save map of moments to quote
		// this is incase theres duplicates due to an extra join from
		// multiple creators (this is a defesnive move against that)
		if currentMomentMap, ok := momentQuoteMap[safeDeref(m.ID)]; ok {
			currentMomentMap[safeDeref(q.ID)] = q
		} else {
			momentQuoteMap[safeDeref(m.ID)] = map[int]*Quote{safeDeref(q.ID): q}
		}
	}

	if !hasRows {
		return nil, ErrNoRecord
	}

	moments := []*Moment{}
	for _, m := range momentMap {
		for _, q := range momentQuoteMap[safeDeref(m.ID)] {
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

	// Get existing quotes for this moment
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

// getExistingQuotes retrieves all quotes for a given moment
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

// insertQuote inserts a new quote
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

// updateQuote updates an existing quote
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

// deleteQuote deletes a quote by ID
func deleteQuote(ctx context.Context, tx pgx.Tx, quoteID *int) error {
	_, err := tx.Exec(ctx, "DELETE FROM quote WHERE id = $1", *quoteID)
	return err
}
