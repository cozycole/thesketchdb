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
	Position      *int
	Actor         *Person
	Character     *Character // if not nil, means character connected to a character's page
	CharacterName *string
	ThumbnailName *string
	ThumbnailFile *multipart.FileHeader
	ProfileFile   *multipart.FileHeader
	// TODO: add attributes about the role
}

type CastModelInterface interface {
	Insert(vidId int, member *CastMember) (int, error)
	InsertThumbnailName(vidId int, name string) error
	GetCastMembers(vidId int) ([]*CastMember, error)
}

type CastModel struct {
	DB *pgxpool.Pool
}

func (m *CastModel) Insert(vidId int, member *CastMember) (int, error) {
	stmt := `
	INSERT INTO cast_members (
	sketch_id, person_id, character_name, character_id, position, img_name)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;`
	result := m.DB.QueryRow(
		context.Background(), stmt, vidId,
		member.Actor.ID, member.CharacterName, member.Character.ID,
		member.Position, member.ThumbnailName)

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

func (m *CastModel) GetCastMembers(vidId int) ([]*CastMember, error) {
	stmt := `
		SELECT p.id, p.slug, p.first, p.last, p.birthdate,
			p.description, p.profile_img, 
			vpr.id, vpr.position, vpr.img_name, vpr.character_name,
			ch.id, ch.slug, ch.name, ch.img_name
		FROM sketch AS v
		JOIN cast_members as vpr
		ON v.id = vpr.sketch_id
		LEFT JOIN person as p
		ON vpr.person_id = p.id
		LEFT JOIN character as ch
		ON vpr.character_id = ch.id
		WHERE v.id = $1
	`
	rows, err := m.DB.Query(context.Background(), stmt, vidId)
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
			&p.ID, &p.Slug, &p.First, &p.Last, &p.BirthDate, &p.Description, &p.ProfileImg,
			&cm.ID, &cm.Position, &cm.ThumbnailName, &cm.CharacterName, &ch.ID, &ch.Slug, &ch.Name, &ch.Image,
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
