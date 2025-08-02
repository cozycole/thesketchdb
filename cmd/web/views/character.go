package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type CharacterPage struct {
	ID            int
	CharacterName string
	Image         string
	PortrayalName string
	PortrayalUrl  string
	Popular       *SketchGallery
}

func CharacterPageView(character *models.Character, popular []*models.Sketch, baseImgUrl string) (*CharacterPage, error) {
	if character.ID == nil {
		return nil, fmt.Errorf("Character ID not defined")
	}

	page := CharacterPage{}

	page.ID = *character.ID

	page.CharacterName = "Missing Character Name"
	if character.Name != nil {
		page.CharacterName = *character.Name
	}

	page.Image = "/static/img/missing-profile.jpg"
	if character.Image != nil {
		page.Image = fmt.Sprintf("%s/character/%s", baseImgUrl, *character.Image)
	}

	if character.Portrayal != nil && character.Portrayal.ID != nil {
		if character.Portrayal.Slug != nil {
			page.PortrayalUrl = fmt.Sprintf(
				"/person/%d/%s",
				*character.Portrayal.ID,
				*character.Portrayal.Slug)
		}

		page.PortrayalName = PrintPersonName(character.Portrayal)
	}
	var err error
	popularPageSize := 12
	page.Popular, err = SketchGalleryView(
		popular,
		baseImgUrl,
		"cast",
		"sub",
		popularPageSize,
	)
	if err != nil {
		return nil, err
	}

	if len(popular) == popularPageSize {
		page.Popular.SeeMore = true
		page.Popular.SeeMoreUrl = fmt.Sprintf(
			"/catalog/sketches?character=%d", safeDeref(character.ID),
		)
	}

	return &page, nil
}

type CharacterGallery struct {
	Cards []*Card
}

func CharacterGalleryView(characters []*models.Character, baseImgUrl string) (*CharacterGallery, error) {
	characterGallery := CharacterGallery{}

	for _, character := range characters {
		characterCard, err := CharacterCardView(character, baseImgUrl)
		if err != nil {
			return nil, err
		}

		characterGallery.Cards = append(characterGallery.Cards, characterCard)
	}

	return &characterGallery, nil
}

func CharacterCardView(character *models.Character, baseImgUrl string) (*Card, error) {
	card := &Card{}

	if character.ID == nil {
		return nil, fmt.Errorf("Character ID not defined")
	}

	if character.Slug == nil {
		return nil, fmt.Errorf("Character slug not defined")
	}

	card.Title = safeDeref(character.Name)

	card.Url = fmt.Sprintf("/character/%d/%s", *character.ID, *character.Slug)
	card.ImageUrl = "/static/img/missing-profile.jpg"
	if character.Image != nil {
		card.ImageUrl = fmt.Sprintf("%s/character/%s", baseImgUrl, *character.Image)
	}

	return card, nil
}
