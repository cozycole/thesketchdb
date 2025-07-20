package main

import (
	"time"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func convertFormToSketch(form *sketchForm) models.Sketch {
	uploadDate, _ := time.Parse(time.DateOnly, form.UploadDate)
	var episode *models.Episode
	var series *models.Series
	var creator *models.Creator
	if form.CreatorID != 0 {
		creator = &models.Creator{
			ID: &form.CreatorID,
		}
	}
	if form.EpisodeID != 0 {
		episode = &models.Episode{ID: &form.EpisodeID}
	}

	if form.SeriesID != 0 {
		series = &models.Series{
			ID: &form.SeriesID,
		}
	}

	return models.Sketch{
		ID:            &form.ID,
		Title:         &form.Title,
		URL:           &form.URL,
		Slug:          &form.Slug,
		ThumbnailFile: form.Thumbnail,
		UploadDate:    &uploadDate,
		Number:        &form.Number,
		Creator:       creator,
		Episode:       episode,
		EpisodeStart:  &form.EpisodeStart,
		Series:        series,
		SeriesPart:    &form.SeriesPart,
	}
}

func convertSketchToForm(sketch *models.Sketch) sketchForm {
	if sketch == nil {
		return sketchForm{}
	}
	var creatorID int
	var creatorName string
	if sketch.Creator != nil {
		creatorID = safeDeref(sketch.Creator.ID)
		creatorName = safeDeref(sketch.Creator.Name)
	}
	var episodeID int
	var episodeName string
	if sketch.Episode != nil {
		episodeID = safeDeref(sketch.Episode.ID)
		episodeName = views.PrintEpisodeName(sketch.Episode)
	}

	var seriesID int
	var seriesName string
	if sketch.Series != nil {
		seriesID = safeDeref(sketch.Series.ID)
		seriesName = safeDeref(sketch.Series.Title)
	}
	return sketchForm{
		ID:           safeDeref(sketch.ID),
		Title:        safeDeref(sketch.Title),
		Slug:         safeDeref(sketch.Slug),
		URL:          safeDeref(sketch.URL),
		UploadDate:   formDate(sketch.UploadDate),
		Number:       safeDeref(sketch.Number),
		CreatorID:    creatorID,
		CreatorInput: creatorName,
		EpisodeID:    episodeID,
		EpisodeInput: episodeName,
		EpisodeStart: safeDeref(sketch.EpisodeStart),
		SeriesID:     seriesID,
		SeriesInput:  seriesName,
		SeriesPart:   safeDeref(sketch.SeriesPart),
	}
}

func convertEpisodeToForm(episode *models.Episode) episodeForm {
	return episodeForm{
		ID:            safeDeref(episode.ID),
		Number:        safeDeref(episode.Number),
		Title:         safeDeref(episode.Title),
		URL:           safeDeref(episode.URL),
		AirDate:       formDate(episode.AirDate),
		ThumbnailName: safeDeref(episode.Thumbnail),
		SeasonId:      safeDeref(episode.SeasonId),
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
		ID:   &form.ID,
		Name: &form.Name,
	}
}

func convertCategoryToForm(category *models.Category) categoryForm {
	return categoryForm{
		ID:   safeDeref(category.ID),
		Name: safeDeref(category.Name),
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
	var categoryId *int
	if form.CategoryID != 0 {
		categoryId = &form.CategoryID
	}
	return models.Tag{
		ID:       &form.ID,
		Name:     &form.Name,
		Category: &models.Category{ID: categoryId},
	}
}

func convertTagtoForm(tag *models.Tag) tagForm {
	var categoryId int
	var categoryName string
	if tag.Category != nil {
		categoryId = safeDeref(tag.Category.ID)
		categoryName = safeDeref(tag.Category.Name)
	}
	return tagForm{
		ID:            safeDeref(tag.ID),
		Name:          safeDeref(tag.Name),
		CategoryID:    categoryId,
		CategoryInput: categoryName,
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
	airDate, err := time.Parse(time.DateOnly, form.AirDate)
	episodeAirDate := &airDate
	if err != nil {
		episodeAirDate = nil
	}
	return models.Episode{
		ID:       &form.ID,
		Title:    &form.Title,
		Number:   &form.Number,
		URL:      &form.URL,
		AirDate:  episodeAirDate,
		SeasonId: &form.SeasonId,
	}
}

func (app *application) convertShowtoForm(show *models.Show) showForm {
	return showForm{
		ID:   safeDeref(show.ID),
		Name: safeDeref(show.Name),
		Slug: safeDeref(show.Slug),
	}
}

func (app *application) convertSeriestoForm(series *models.Series) seriesForm {
	return seriesForm{
		ID:            safeDeref(series.ID),
		Title:         safeDeref(series.Title),
		ThumbnailName: safeDeref(series.ThumbnailName),
	}
}

func (app *application) convertFormtoSeries(form *seriesForm) models.Series {
	return models.Series{
		ID:            &form.ID,
		Title:         &form.Title,
		ThumbnailName: &form.ThumbnailName,
	}

}
