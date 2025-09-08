package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type PersonPage struct {
	ID             int
	Name           string
	Image          string
	BirthDate      string
	Age            int
	Professions    string
	SketchCount    int
	PortrayalCount int
	CharacterCount int
	CreatorCount   int
	Popular        *SketchGallery
	ShowCreators   ProfileGallery
}

func PersonPageView(
	person *models.Person,
	stats *models.PersonStats,
	popular []*models.Sketch,
	creatorShowCounts []*models.CreatorShowCounts,
	baseImgUrl string) (*PersonPage, error) {
	page := PersonPage{}

	if person.ID == nil {
		return nil, fmt.Errorf("Person ID not defined")
	}

	page.ID = *person.ID
	page.Name = PrintPersonName(person)

	if person.ProfileImg != nil {
		page.Image = fmt.Sprintf("%s/person/medium/%s", baseImgUrl, *person.ProfileImg)
	}

	if person.BirthDate != nil {
		page.BirthDate = humanDate(person.BirthDate)
		page.Age = getAge(person.BirthDate)
	}

	if person.Professions != nil {
		page.Professions = *person.Professions
	}

	page.SketchCount = stats.SketchCount
	page.PortrayalCount = stats.PortrayalCount
	page.CharacterCount = stats.CharacterCount
	page.CreatorCount = stats.CreatorCount

	var err error
	page.Popular, err = SketchGalleryView(
		popular,
		baseImgUrl,
		"cast",
		"sub",
		12,
	)
	if len(popular) == 12 {
		page.Popular.SeeMore = true
		page.Popular.SeeMoreUrl = fmt.Sprintf("/catalog/sketches?person=%d", page.ID)
	}
	if err != nil {
		return nil, err
	}

	page.ShowCreators = ShowCreatorGalleryView(
		creatorShowCounts,
		safeDeref(person.ID),
		baseImgUrl,
	)

	return &page, nil

}

type PersonGallery struct {
	Cards []*Card
}

func PersonGalleryView(people []*models.Person, baseImgUrl string) (*PersonGallery, error) {
	personGallery := PersonGallery{}

	for _, person := range people {
		personCard, err := PersonCardView(person, baseImgUrl)
		if err != nil {
			return nil, err
		}

		personGallery.Cards = append(personGallery.Cards, personCard)
	}

	return &personGallery, nil
}

func PersonCardView(person *models.Person, baseImgUrl string) (*Card, error) {
	card := &Card{}

	if person.ID == nil {
		return nil, fmt.Errorf("Person ID not defined")
	}

	if person.Slug == nil {
		return nil, fmt.Errorf("Person slug not defined")
	}

	card.Title = PrintPersonName(person)
	card.Url = fmt.Sprintf("/person/%d/%s", *person.ID, *person.Slug)
	card.ImageUrl = "/static/img/missing-profile.jpg"
	if person.ProfileImg != nil {
		card.ImageUrl = fmt.Sprintf("%s/person/medium/%s", baseImgUrl, *person.ProfileImg)
	}

	return card, nil
}

type ProfileGallery struct {
	Cards []*Card
}

type Card struct {
	Url      string
	ImageUrl string
	Title    string
	Subtitle string
}

func ShowCreatorGalleryView(creatorShows []*models.CreatorShowCounts, personId int, baseImageUrl string) ProfileGallery {
	cards := []*Card{}
	for _, cs := range creatorShows {
		subtitle := fmt.Sprintf(
			"Featured in %s", sketchCountLabel(safeDeref(cs.Count)),
		)
		url := fmt.Sprintf("/catalog/sketches?person=%d&%s=%d",
			personId,
			safeDeref(cs.Type),
			safeDeref(cs.ID),
		)

		imageUrl := fmt.Sprintf("%s/%s/medium/%s",
			baseImageUrl,
			safeDeref(cs.Type),
			safeDeref(cs.ImageName),
		)

		card := Card{
			Url:      url,
			ImageUrl: imageUrl,
			Title:    safeDeref(cs.Name),
			Subtitle: subtitle,
		}

		cards = append(cards, &card)
	}

	return ProfileGallery{Cards: cards}
}
