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
	// 1) Cast Image 2) Character Profile 3) Actor Profile 4) Missing Profile
	card.Image = fmt.Sprintf("/static/img/missing-profile.jpg")
	if member.Actor != nil {
		if member.Actor.ID != nil && member.Actor.Slug != nil {
			card.ActorUrl = fmt.Sprintf(
				"/person/%d/%s",
				*member.Actor.ID,
				*member.Actor.Slug)

			if member.Actor.ProfileImg != nil {
				card.Image = fmt.Sprintf(
					"%s/person/%s",
					baseImgUrl,
					*member.Actor.ProfileImg)
			}
		}
	}

	card.ActorName = printPersonName(member.Actor)

	if member.Character != nil {
		if member.Character.ID != nil && member.Character.Slug != nil {
			card.CharacterUrl = fmt.Sprintf(
				"/character/%d/%s",
				*member.Character.ID,
				*member.Character.Slug)

			if member.Character.Image != nil {
				card.Image = fmt.Sprintf(
					"%s/character/%s",
					baseImgUrl,
					*member.Character.Image)
			}
		}
	}

	if member.ThumbnailName != nil {
		card.Image =
			fmt.Sprintf(
				"%s/cast/profile/%s",
				baseImgUrl,
				*member.ThumbnailName)

	}

	card.CardUrl = card.ActorUrl
	if card.CharacterUrl != "" {
		card.CardUrl = card.CharacterUrl
	}

	return &card, nil
}
