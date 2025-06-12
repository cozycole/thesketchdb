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
}

func CreatorPageView(creator *models.Creator, popular []*models.Video, baseImgUrl string) (*CreatorPage, error) {
	if creator.ID == nil {
		return nil, fmt.Errorf("Creator has no defined ID")
	}

	page := CreatorPage{}
	page.ID = *creator.ID

	page.CreatorName = safeDerefString(creator.Name)
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

	return &page, nil
}

type CreatorGallery struct {
	CreatorCards []*CreatorCard
}

type CreatorCard struct {
	Name  string
	Url   string
	Image string
}

func CreatorGalleryView(creators []*models.Creator, baseImgUrl string) (*CreatorGallery, error) {
	creatorGallery := CreatorGallery{}

	for _, creator := range creators {
		creatorCard, err := CreatorCardView(creator, baseImgUrl)
		if err != nil {
			return nil, err
		}

		creatorGallery.CreatorCards = append(creatorGallery.CreatorCards, creatorCard)
	}

	return &creatorGallery, nil
}

func CreatorCardView(creator *models.Creator, baseImgUrl string) (*CreatorCard, error) {
	card := &CreatorCard{}

	if creator.ID == nil {
		return nil, fmt.Errorf("Creator ID not defined")
	}

	if creator.Slug == nil {
		return nil, fmt.Errorf("Creator slug not defined")
	}

	card.Name = "Missing Creator Name"
	if creator.Name != nil {
		card.Name = *creator.Name
	}

	card.Url = fmt.Sprintf("/creator/%d/%s", *creator.ID, *creator.Slug)
	card.Image = "/static/img/missing-profile.jpg"
	if creator.ProfileImage != nil {
		card.Image = fmt.Sprintf("%s/creator/%s", baseImgUrl, *creator.ProfileImage)
	}

	return card, nil
}
