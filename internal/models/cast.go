package models

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CastMember struct {
	ID            *int
	SketchID      *int
	Position      *int
	Actor         *Person
	Character     *Character // if not nil, means character connected to a character's page
	CharacterName *string
	CastRole      *string
	MinorRole     *bool
	ThumbnailName *string
	ProfileImg    *string
	ThumbnailFile *multipart.FileHeader
	ProfileFile   *multipart.FileHeader
	Tags          []*Tag
}

type CastModelInterface interface {
	Delete(id int) error
	GetById(id int) (*CastMember, error)
	Insert(sketchId int, member *CastMember) (int, error)
	InsertThumbnailName(sketchId int, name string) error
	GetCastMembers(sketchId int) ([]*CastMember, error)
	BatchUpdateCastTags(castID int, tags []*Tag) error
	Update(member *CastMember) error
	UpdatePositions([]int) error
}

type CastModel struct {
	DB *pgxpool.Pool
}

func (m *CastModel) Delete(id int) error {
	stmt := `
	DELETE FROM cast_members WHERE id = $1
	`
	_, err := m.DB.Exec(context.Background(), stmt, id)
	return err
}
func (m *CastModel) GetById(id int) (*CastMember, error) {
	stmt := `
	SELECT cm.id, cm.sketch_id, cm.position, cm.character_name, cm.role, cm.minor,
	cm.thumbnail_name, cm.profile_img,
	p.id, p.first, p.last, p.profile_img,
	ch.id, ch.name, ch.img_name,
	t.id, t.name, t.slug, t.type,
	c.id, c.name, c.slug
	FROM cast_members as cm
	LEFT JOIN person as p ON cm.person_id = p.id
	LEFT JOIN character as ch ON cm.character_id = ch.id
	LEFT JOIN cast_tags_rel as ctr ON cm.id = ctr.cast_id
	LEFT JOIN tags as t ON ctr.tag_id = t.id
	LEFT JOIN categories as c ON t.category_id = c.id
	WHERE cm.id = $1
	`

	c := CastMember{}
	p := Person{}
	ch := Character{}
	tags := []*Tag{}
	rows, err := m.DB.Query(context.Background(), stmt, id)
	for rows.Next() {
		t := Tag{}
		ca := Category{}
		rows.Scan(
			&c.ID, &c.SketchID, &c.Position, &c.CharacterName, &c.CastRole,
			&c.MinorRole, &c.ThumbnailName, &c.ProfileImg, &p.ID, &p.First,
			&p.Last, &p.ProfileImg, &ch.ID, &ch.Name, &ch.Image,
			&t.ID, &t.Name, &t.Slug, &t.Type,
			&ca.ID, &ca.Name, &ca.Slug,
		)

		if t.ID != nil {
			t.Category = &ca
			tags = append(tags, &t)
		}
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	if p.ID != nil {
		c.Actor = &p
	}

	if ch.ID != nil {
		c.Character = &ch
	}

	c.Tags = tags

	return &c, nil
}

func (m *CastModel) Insert(sketchId int, member *CastMember) (int, error) {
	var actorId, characterId *int
	if member.Actor != nil && safeDeref(member.Actor.ID) != 0 {
		actorId = member.Actor.ID
	}

	if member.Character != nil && safeDeref(member.Character.ID) != 0 {
		characterId = member.Character.ID
	}

	stmt := `
	INSERT INTO cast_members (
	sketch_id, person_id, character_name, character_id, position, role, minor,
	thumbnail_name, profile_img)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id;`
	result := m.DB.QueryRow(
		context.Background(), stmt, sketchId, actorId, member.CharacterName,
		characterId, member.Position, member.CastRole, member.MinorRole,
		member.ThumbnailName, member.ProfileImg)

	var id int
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *CastModel) InsertThumbnailName(castId int, name string) error {
	stmt := `UPDATE cast_members SET img_name = $1 WHERE id = $2`
	_, err := m.DB.Exec(context.Background(), stmt, name, castId)
	return err
}

func (m *CastModel) GetCastMembers(sketchId int) ([]*CastMember, error) {
	stmt := `
		SELECT p.id, p.slug, p.first, p.last, p.birthdate,
		p.description, p.profile_img, 
		cm.id, cm.position, cm.thumbnail_name, cm.profile_img, 
		cm.character_name, cm.role, cm.minor,
		ch.id, ch.slug, ch.name, ch.img_name,
		t.id, t.name, t.slug, t.type,
		c.id, c.name, c.slug
		FROM sketch AS v
		JOIN cast_members as cm ON v.id = cm.sketch_id
		LEFT JOIN person as p ON cm.person_id = p.id
		LEFT JOIN character as ch ON cm.character_id = ch.id
		LEFT JOIN cast_tags_rel as ctr ON cm.id = ctr.cast_id
		LEFT JOIN tags as t ON ctr.tag_id = t.id
		LEFT JOIN categories as c ON t.category_id = c.id
		WHERE v.id = $1
		ORDER BY cm.position asc
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

	memberTagMap := map[int]map[int]*Tag{}
	memberMap := map[int]*CastMember{}
	for rows.Next() {
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		t := Tag{}
		c := Category{}
		err := rows.Scan(
			&p.ID, &p.Slug, &p.First, &p.Last, &p.BirthDate,
			&p.Description, &p.ProfileImg, &cm.ID, &cm.Position,
			&cm.ThumbnailName, &cm.ProfileImg, &cm.CharacterName,
			&cm.CastRole, &cm.MinorRole, &ch.ID, &ch.Slug, &ch.Name, &ch.Image,
			&t.ID, &t.Name, &t.Slug, &t.Type,
			&c.ID, &c.Name, &c.Slug,
		)
		if err != nil {
			return nil, err
		}

		if cm.ID == nil {
			continue
		}

		cm.Actor = p
		cm.Character = ch

		memberMap[*cm.ID] = cm

		if t.ID == nil {
			continue
		}

		if tagMap, ok := memberTagMap[*cm.ID]; ok {
			tagMap[*t.ID] = &t
		} else {
			memberTagMap[*cm.ID] = map[int]*Tag{*t.ID: &t}
		}
	}

	members := []*CastMember{}
	for cast_id, cm := range memberMap {
		if tagMap, ok := memberTagMap[cast_id]; ok {
			tags := []*Tag{}
			for _, tag := range tagMap {
				tags = append(tags, tag)
			}
			cm.Tags = tags
		}
		members = append(members, cm)
	}
	sort.Slice(members, func(i, j int) bool {
		return safeDeref(members[i].Position) < safeDeref(members[j].Position)
	})

	return members, nil
}

func (m *CastModel) Update(member *CastMember) error {
	stmt := `
	UPDATE cast_members SET person_id = $1, character_name = $2, character_id = $3,
	role = $4, thumbnail_name = $5, profile_img = $6, minor = $7
	WHERE id = $8
	`
	var personId *int
	if member.Actor != nil && safeDeref(member.Actor.ID) != 0 {
		personId = member.Actor.ID

	}
	_, err := m.DB.Exec(context.Background(), stmt, personId, member.CharacterName,
		member.Character.ID, member.CastRole, member.ThumbnailName, member.ProfileImg,
		member.MinorRole, member.ID,
	)
	return err
}

func (m *CastModel) UpdatePositions(castIds []int) error {
	stmt := `
	UPDATE cast_members AS cm
	SET position = data.pos
	FROM (
		SELECT * FROM unnest($1::int[]) WITH ORDINALITY
	) AS data(id, pos)
	WHERE cm.id = data.id;
	`
	_, err := m.DB.Exec(context.Background(), stmt, castIds)

	return err
}

func (m *CastModel) BatchUpdateCastTags(castID int, tags []*Tag) error {
	ctx := context.Background()
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get existing tag associations for this cast
	existingTagIDs, err := getExistingCastTags(ctx, tx, castID)
	if err != nil {
		return fmt.Errorf("failed to get existing cast tags: %w", err)
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
		err = insertCastTagAssociations(ctx, tx, castID, tagsToInsert)
		if err != nil {
			return fmt.Errorf("failed to insert cast tag associations: %w", err)
		}
	}

	// Delete removed tag associations
	if len(tagsToDelete) > 0 {
		err = deleteCastTagAssociations(ctx, tx, castID, tagsToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete cast tag associations: %w", err)
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getExistingCastTags(ctx context.Context, tx pgx.Tx, castID int) ([]int, error) {
	rows, err := tx.Query(ctx, `
		SELECT tag_id 
		FROM cast_tags_rel 
		WHERE cast_id = $1 
		ORDER BY tag_id`, castID)
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

func insertCastTagAssociations(ctx context.Context, tx pgx.Tx, castID int, tagIDs []int) error {
	// Use batch insert for efficiency
	batch := &pgx.Batch{}

	for _, tagID := range tagIDs {
		batch.Queue("INSERT INTO cast_tags_rel (cast_id, tag_id) VALUES ($1, $2)", castID, tagID)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	// Execute all batched queries
	for i := range len(tagIDs) {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert cast_tag association for tag %d: %w", tagIDs[i], err)
		}
	}

	return nil
}

func deleteCastTagAssociations(ctx context.Context, tx pgx.Tx, castID int, tagIDs []int) error {
	// Use a single query with IN clause for efficiency
	query := `DELETE FROM cast_tags_rel WHERE cast_id = $1 AND tag_id = ANY($2)`

	_, err := tx.Exec(ctx, query, castID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to delete cast tag associations: %w", err)
	}

	return nil
}
