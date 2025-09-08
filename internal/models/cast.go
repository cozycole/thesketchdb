package models

import (
	"context"
	"errors"
	"mime/multipart"

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
}

type CastModelInterface interface {
	Delete(id int) error
	GetById(id int) (*CastMember, error)
	Insert(sketchId int, member *CastMember) (int, error)
	InsertThumbnailName(sketchId int, name string) error
	GetCastMembers(sketchId int) ([]*CastMember, error)
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
	SELECT c.id, c.sketch_id, c.position, c.character_name, c.role, c.minor,
	c.thumbnail_name, c.profile_img,
	p.id, p.first, p.last, p.profile_img,
	ch.id, ch.name, ch.img_name
	FROM cast_members as c
	LEFT JOIN person as p
	ON c.person_id = p.id
	LEFT JOIN character as ch 
	ON c.character_id = ch.id
	WHERE c.id = $1
	`
	c := CastMember{}
	p := Person{}
	ch := Character{}
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(
		&c.ID, &c.SketchID, &c.Position, &c.CharacterName, &c.CastRole,
		&c.MinorRole, &c.ThumbnailName, &c.ProfileImg, &p.ID, &p.First,
		&p.Last, &p.ProfileImg, &ch.ID, &ch.Name, &ch.Image,
	)

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
			ch.id, ch.slug, ch.name, ch.img_name
		FROM sketch AS v
		JOIN cast_members as cm
		ON v.id = cm.sketch_id
		LEFT JOIN person as p
		ON cm.person_id = p.id
		LEFT JOIN character as ch
		ON cm.character_id = ch.id
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

	members := []*CastMember{}
	for rows.Next() {
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		err := rows.Scan(
			&p.ID, &p.Slug, &p.First, &p.Last, &p.BirthDate,
			&p.Description, &p.ProfileImg, &cm.ID, &cm.Position,
			&cm.ThumbnailName, &cm.ProfileImg, &cm.CharacterName,
			&cm.CastRole, &cm.MinorRole, &ch.ID, &ch.Slug, &ch.Name, &ch.Image,
		)
		if err != nil {
			return nil, err
		}

		cm.Actor = p
		cm.Character = ch
		members = append(members, cm)
	}

	return members, nil
}

func (m *CastModel) Update(member *CastMember) error {
	stmt := `
	UPDATE cast_members SET person_id = $1, character_name = $2, character_id = $3,
	role = $4, thumbnail_name = $5, profile_img = $6, minor = $7
	WHERE id = $8
	`
	_, err := m.DB.Exec(context.Background(), stmt, member.Actor.ID, member.CharacterName,
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
