package main

import (
	"time"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func convertFormToSketch(form *sketchForm) models.Sketch {
	uploadDate, _ := time.Parse(time.DateOnly, form.UploadDate)
	return models.Sketch{
		ID:            &form.ID,
		Title:         &form.Title,
		URL:           &form.URL,
		Slug:          &form.Slug,
		ThumbnailFile: form.Thumbnail,
		Rating:        &form.Rating,
		UploadDate:    &uploadDate,
		Creator: &models.Creator{
			ID: &form.CreatorID,
		},
	}
}

func convertSketchToForm(sketch *models.Sketch) sketchForm {
	var creatorID int
	var creatorName string
	if sketch.Creator != nil {
		creatorID = safeDeref(sketch.Creator.ID)
		creatorName = safeDeref(sketch.Creator.Name)
	}
	return sketchForm{
		ID:           safeDeref(sketch.ID),
		Title:        safeDeref(sketch.Title),
		Slug:         safeDeref(sketch.Slug),
		URL:          safeDeref(sketch.URL),
		UploadDate:   formDate(sketch.UploadDate),
		CreatorID:    creatorID,
		CreatorInput: creatorName,
	}
}

func convertFormtoCastMember(form *castForm) models.CastMember {
	actor := models.Person{ID: &form.PersonID}
	character := models.Character{}
	if form.CharacterID != 0 {
		character.ID = &form.CharacterID
	}
	return models.CastMember{
		ID:            &form.ID,
		Actor:         &actor,
		Character:     &character,
		CastRole:      &form.CastRole,
		MinorRole:     &form.MinorRole,
		CharacterName: &form.CharacterName,
		ThumbnailFile: form.CharacterThumbnail,
		ProfileFile:   form.CharacterProfile,
	}
}

func convertCastMembertoForm(member *models.CastMember) castForm {
	var personID, characterID int
	var personName, characterName string
	if member.Actor != nil {
		personID = safeDeref(member.Actor.ID)
		personName = views.PrintPersonName(member.Actor)
	}

	if member.Character != nil {
		characterID = safeDeref(member.Character.ID)
		characterName = safeDeref(member.Character.Name)
	}

	return castForm{
		ID:             safeDeref(member.ID),
		PersonID:       personID,
		PersonInput:    personName,
		CastRole:       safeDeref(member.CastRole),
		MinorRole:      safeDeref(member.MinorRole),
		CharacterName:  safeDeref(member.CharacterName),
		CharacterID:    characterID,
		CharacterInput: characterName,
		ThumbnailName:  safeDeref(member.ThumbnailName),
		ProfileImage:   safeDeref(member.ProfileImg),
	}
}

func convertFormtoCategory(form *categoryForm) models.Category {
	return models.Category{
		Name: &form.Name,
	}
}

func convertFormtoPerson(form *personForm) models.Person {
	var id *int
	if form.ID != 0 {
		id = &form.ID
	}
	bdate, _ := time.Parse(time.DateOnly, form.BirthDate)
	return models.Person{
		ID:          id,
		First:       &form.First,
		Last:        &form.Last,
		BirthDate:   &bdate,
		Professions: &form.Professions,
	}
}

func convertPersontoForm(person *models.Person) personForm {
	return personForm{
		ID:          safeDeref(person.ID),
		First:       safeDeref(person.First),
		Last:        safeDeref(person.Last),
		Professions: safeDeref(person.Professions),
		BirthDate:   formDate(person.BirthDate),
		ImageUrl:    safeDeref(person.ProfileImg),
	}
}

func convertFormtoCharacter(form *characterForm) models.Character {
	var id *int
	if form.ID != 0 {
		id = &form.ID
	}
	return models.Character{
		ID:        id,
		Name:      &form.Name,
		Type:      &form.Type,
		Portrayal: &models.Person{ID: &form.PersonID},
	}
}

func convertCharactertoForm(character *models.Character) characterForm {
	var personInput string
	var personId int
	if character.Portrayal != nil {
		if character.Portrayal.ID != nil {
			personId = *character.Portrayal.ID
		}

		personInput = views.PrintPersonName(character.Portrayal)
	}

	return characterForm{
		ID:          safeDeref(character.ID),
		Name:        safeDeref(character.Name),
		Type:        safeDeref(character.Type),
		ImageUrl:    safeDeref(character.Image),
		PersonID:    personId,
		PersonInput: personInput,
	}
}

func convertFormtoCreator(form *creatorForm) models.Creator {
	var id *int
	if form.ID != 0 {
		id = &form.ID
	}
	date, _ := time.Parse(time.DateOnly, form.EstablishedDate)
	return models.Creator{
		ID:              id,
		Name:            &form.Name,
		URL:             &form.URL,
		EstablishedDate: &date,
	}
}

func convertCreatortoForm(creator *models.Creator) creatorForm {
	return creatorForm{
		ID:              safeDeref(creator.ID),
		Name:            safeDeref(creator.Name),
		URL:             safeDeref(creator.URL),
		EstablishedDate: formDate(creator.EstablishedDate),
		ImageUrl:        safeDeref(creator.ProfileImage),
	}
}

func convertFormtoTag(form *tagForm) models.Tag {
	var categoryId int
	if form.CategoryId != 0 {
		categoryId = form.CategoryId
	}
	return models.Tag{
		Name:     &form.Name,
		Category: &models.Category{ID: &categoryId},
	}
}

func (app *application) convertFormtoSketchTags(form *sketchTagsForm) ([]*models.Tag, error) {
	var tags []*models.Tag
	for _, tagId := range form.TagIds {
		tag, err := app.tags.Get(tagId)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (app *application) convertFormtoShow(form *showForm) models.Show {
	return models.Show{
		Name: &form.Name,
		Slug: &form.Slug,
	}

}

func (app *application) convertFormtoEpisode(form *episodeForm) models.Episode {
	airDate, _ := time.Parse(time.DateOnly, form.AirDate)
	return models.Episode{
		ID:       &form.ID,
		Title:    &form.Title,
		Number:   &form.Number,
		AirDate:  &airDate,
		SeasonId: &form.SeasonId,
	}
}
