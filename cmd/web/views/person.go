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
	SketchCount    int
	PortrayalCount int
	CharacterCount int
	CreatorCount   int
	Popular        *SketchGallery
}

func PersonPageView(
	person *models.Person,
	stats *models.PersonStats,
	popular []*models.Sketch,
	baseImgUrl string) (*PersonPage, error) {
	page := PersonPage{}

	if person.ID == nil {
		return nil, fmt.Errorf("Person ID not defined")
	}

	page.ID = *person.ID
	page.Name = printPersonName(person)

	if person.ProfileImg != nil {
		page.Image = fmt.Sprintf("%s/person/%s", baseImgUrl, *person.ProfileImg)
	}

	if person.BirthDate != nil {
		page.BirthDate = humanDate(person.BirthDate)
		page.Age = getAge(person.BirthDate)
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
	if err != nil {
		return nil, err
	}

	return &page, nil

}

type PersonGallery struct {
	PersonCards []*PersonCard
}

type PersonCard struct {
	Name  string
	Url   string
	Image string
}

func PersonGalleryView(people []*models.Person, baseImgUrl string) (*PersonGallery, error) {
	personGallery := PersonGallery{}

	for _, person := range people {
		personCard, err := PersonCardView(person, baseImgUrl)
		if err != nil {
			return nil, err
		}

		personGallery.PersonCards = append(personGallery.PersonCards, personCard)
	}

	return &personGallery, nil
}

func PersonCardView(person *models.Person, baseImgUrl string) (*PersonCard, error) {
	card := &PersonCard{}

	if person.ID == nil {
		return nil, fmt.Errorf("Person ID not defined")
	}

	if person.Slug == nil {
		return nil, fmt.Errorf("Person slug not defined")
	}

	card.Name = printPersonName(person)
	if card.Name == "" {
		card.Name = "Missing Name"
	}

	card.Url = fmt.Sprintf("/person/%d/%s", *person.ID, *person.Slug)
	card.Image = "/static/img/missing-profile.jpg"
	if person.ProfileImg != nil {
		card.Image = fmt.Sprintf("%s/person/%s", baseImgUrl, *person.ProfileImg)
	}

	return card, nil
}
