package views

import (
	"fmt"
	"html/template"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/services/moviedb"
	"sketchdb.cozycole.net/internal/services/wikipedia"
)

type PersonPage struct {
	ID                    int
	Name                  string
	Image                 string
	BirthDate             string
	Description           template.HTML
	Age                   int
	WikiLink              string
	IMDbUrl               string
	Professions           string
	SketchCount           string
	PortrayalCount        string
	ImpressionCount       string
	DisplayCharacterStats bool
	OriginalCount         string
	CharacterCount        string
	CreatorCount          string
	ShowCount             string
	Popular               *SketchGallery
	ShowCreators          ProfileGallery
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

	if wikiPage := safeDeref(person.WikiPage); wikiPage != "" {
		page.WikiLink = fmt.Sprintf(wikipedia.URL_TEMPLATE, wikiPage)
	}

	page.Professions = safeDeref(person.Professions)
	page.Description = template.HTML(safeDeref(person.Description))

	page.SketchCount = sketchCountLabel(stats.SketchCount)

	if stats.CreatorCount != 0 {
		page.CreatorCount = countLabel("Creator", stats.CreatorCount)
	}

	if stats.ShowCount != 0 {
		page.ShowCount = countLabel("Show", stats.ShowCount)
	}

	if safeDeref(person.IMDbID) != "" {
		page.IMDbUrl = moviedb.BuildIMDbURL(*person.IMDbID)
	}

	if stats.PortrayalCount != 0 || stats.ImpressionCount != 0 || stats.OriginalCount != 0 {
		page.DisplayCharacterStats = true
	}

	if stats.ImpressionCount != 0 {
		page.ImpressionCount = countLabel("Impression", stats.ImpressionCount)
	}

	if stats.OriginalCount != 0 {
		page.OriginalCount = countLabel("Original", stats.OriginalCount)
	}

	if stats.PortrayalCount != 0 {
		page.PortrayalCount = countLabel("Potrayal", stats.PortrayalCount)
	}

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
