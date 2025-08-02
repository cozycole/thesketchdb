package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type CreatorPage struct {
	ID              int
	CreatorName     string
	Image           string
	EstablishedDate string
	Popular         *SketchGallery
	CastSection     *PersonGallery
}

func CreatorPageView(
	creator *models.Creator,
	popular []*models.Sketch,
	cast []*models.Person,
	baseImgUrl string) (*CreatorPage, error) {
	if creator.ID == nil {
		return nil, fmt.Errorf("Creator has no defined ID")
	}

	page := CreatorPage{}
	page.ID = *creator.ID

	page.CreatorName = safeDeref(creator.Name)
	if page.CreatorName == "" {
		page.CreatorName = "Missing Creator Name"
	}

	page.Image = "/static/img/missing-profile.jpg"
	if creator.ProfileImage != nil {
		page.Image = fmt.Sprintf(
			"%s/creator/%s",
			baseImgUrl,
			*creator.ProfileImage,
		)
	}

	page.EstablishedDate = humanDate(creator.EstablishedDate)

	var err error
	page.Popular, err = SketchGalleryView(
		popular,
		baseImgUrl,
		"base",
		"sub",
		12,
	)
	if err != nil {
		return nil, err
	}

	page.CastSection, err = PersonGalleryView(cast, baseImgUrl)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

type CreatorGallery struct {
	Cards []*Card
}

func CreatorGalleryView(creators []*models.Creator, baseImgUrl string) (*CreatorGallery, error) {
	creatorGallery := CreatorGallery{}

	for _, creator := range creators {
		creatorCard, err := CreatorCardView(creator, baseImgUrl)
		if err != nil {
			return nil, err
		}

		creatorGallery.Cards = append(creatorGallery.Cards, creatorCard)
	}

	return &creatorGallery, nil
}

func CreatorCardView(creator *models.Creator, baseImgUrl string) (*Card, error) {
	card := &Card{}

	if creator.ID == nil {
		return nil, fmt.Errorf("Creator ID not defined")
	}

	if creator.Slug == nil {
		return nil, fmt.Errorf("Creator slug not defined")
	}

	card.Title = safeDeref(creator.Name)

	card.Url = fmt.Sprintf("/creator/%d/%s", safeDeref(creator.ID), safeDeref(creator.Slug))
	card.ImageUrl = "/static/img/missing-profile.jpg"
	if creator.ProfileImage != nil {
		card.ImageUrl = fmt.Sprintf("%s/creator/%s", baseImgUrl, *creator.ProfileImage)
	}

	return card, nil
}
