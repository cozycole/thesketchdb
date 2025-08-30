package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type CastGallery struct {
	CastCards []*CastCard
}

type CastCard struct {
	CardUrl       string
	ActorName     string
	ActorUrl      string
	CharacterName string
	CharacterUrl  string
	CastRole      string
	Image         string
}

func CastGalleryView(cast []*models.CastMember, baseImgUrl string) (CastGallery, error) {
	gallery := CastGallery{}
	for _, c := range cast {
		card, err := CastCardView(c, baseImgUrl)
		if err != nil {
			return gallery, err
		}

		gallery.CastCards = append(gallery.CastCards, card)
	}

	return gallery, nil
}

func CastCardView(member *models.CastMember, baseImgUrl string) (*CastCard, error) {
	card := CastCard{}

	if member.CharacterName != nil {
		card.CharacterName = *member.CharacterName
	}

	// Card Image can be, based on priority:
	// 1) Cast Image 2) Actor Profile 3) Missing Profile
	card.Image = fmt.Sprintf("/static/img/missing-profile.jpg")
	if member.Actor != nil {
		if member.Actor.ID != nil && member.Actor.Slug != nil {
			card.ActorUrl = fmt.Sprintf(
				"/person/%d/%s",
				*member.Actor.ID,
				*member.Actor.Slug)

			if member.Actor.ProfileImg != nil {
				card.Image = fmt.Sprintf(
					"%s/person/small/%s",
					baseImgUrl,
					*member.Actor.ProfileImg)
			}
		}
	}

	card.ActorName = PrintPersonName(member.Actor)
	card.CastRole = uppercaseFirst(safeDeref(member.CastRole))

	var characterType string
	if member.Character != nil {
		if member.Character.ID != nil && member.Character.Slug != nil {
			card.CharacterUrl = fmt.Sprintf(
				"/character/%d/%s",
				*member.Character.ID,
				*member.Character.Slug)
		}
		characterType = safeDeref(member.Character.Type)
	}

	if member.ProfileImg != nil {
		card.Image =
			fmt.Sprintf(
				"%s/cast/profile/medium/%s",
				baseImgUrl,
				*member.ProfileImg)
	}

	card.CardUrl = card.ActorUrl

	// we only want to route to a character page on card click if it's
	// an interesting character (i.e. not generic)
	useCharacterUrl := characterType != "" && characterType != "generic"
	if card.CharacterUrl != "" && useCharacterUrl {
		card.CardUrl = card.CharacterUrl
	}

	return &card, nil
}

type CastTable struct {
	SketchID int
	CastRows []CastRow
}

type CastRow struct {
	ID                int
	ImageUrl          string
	CharacterName     string
	PersonName        string
	PersonUrl         string
	CharacterPageName string
	CharacterUrl      string
	CastRole          string
	MinorRole         bool
}

func CastTableView(cast []*models.CastMember, sketchID int, baseImgUrl string) CastTable {
	castTable := CastTable{}
	castTable.SketchID = sketchID
	for _, c := range cast {
		row := CastRow{}
		row.ID = safeDeref(c.ID)
		if c.ProfileImg == nil && c.Actor != nil {
			row.ImageUrl = fmt.Sprintf("%s/person/small/%s", baseImgUrl, safeDeref(c.Actor.ProfileImg))
		} else {
			row.ImageUrl = fmt.Sprintf("%s/cast/profile/small/%s", baseImgUrl, safeDeref(c.ProfileImg))
		}

		row.CharacterName = safeDeref(c.CharacterName)
		if c.Actor != nil {
			row.PersonName = PrintPersonName(c.Actor)
			row.PersonUrl = fmt.Sprintf(
				"/person/%d/%s",
				safeDeref(c.Actor.ID),
				safeDeref(c.Actor.Slug),
			)
		}

		if c.Character != nil {
			row.CharacterPageName = safeDeref(c.Character.Name)
			row.CharacterUrl = fmt.Sprintf(
				"/character/%d/%s",
				safeDeref(c.Character.ID),
				safeDeref(c.Character.Slug),
			)
		}

		row.CastRole = uppercaseFirst(safeDeref(c.CastRole))
		row.MinorRole = safeDeref(c.MinorRole)

		castTable.CastRows = append(castTable.CastRows, row)
	}

	return castTable
}
