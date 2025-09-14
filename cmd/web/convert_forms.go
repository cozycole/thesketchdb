package main

import (
	"fmt"
	"time"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func convertFormToSketch(form *sketchForm) models.Sketch {
	uploadDate, _ := time.Parse(time.DateOnly, form.UploadDate)
	var creator *models.Creator
	if form.CreatorID != 0 {
		creator = &models.Creator{
			ID: &form.CreatorID,
		}
	}

	var episode *models.Episode
	if form.EpisodeID != 0 {
		episode = &models.Episode{ID: &form.EpisodeID}
	}

	var series *models.Series
	if form.SeriesID != 0 {
		series = &models.Series{
			ID: &form.SeriesID,
		}
	}

	var recurring *models.Recurring
	if form.RecurringID != 0 {
		recurring = &models.Recurring{
			ID: &form.RecurringID,
		}
	}

	return models.Sketch{
		ID:            &form.ID,
		Title:         &form.Title,
		URL:           &form.URL,
		Slug:          &form.Slug,
		Duration:      &form.Duration,
		Description:   &form.Description,
		Transcript:    &form.Transcript,
		Diarization:   &form.Diarization,
		ThumbnailFile: form.Thumbnail,
		UploadDate:    &uploadDate,
		Number:        &form.Number,
		Popularity:    &form.Popularity,
		Creator:       creator,
		Episode:       episode,
		EpisodeStart:  &form.EpisodeStart,
		Series:        series,
		SeriesPart:    &form.SeriesPart,
		Recurring:     recurring,
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

	var recurringID int
	var recurringName string
	if sketch.Recurring != nil {
		recurringID = safeDeref(sketch.Recurring.ID)
		recurringName = safeDeref(sketch.Recurring.Title)
	}
	return sketchForm{
		ID:             safeDeref(sketch.ID),
		Slug:           safeDeref(sketch.Slug),
		Title:          safeDeref(sketch.Title),
		URL:            safeDeref(sketch.URL),
		Duration:       safeDeref(sketch.Duration),
		Description:    safeDeref(sketch.Description),
		Transcript:     safeDeref(sketch.Transcript),
		Diarization:    safeDeref(sketch.Diarization),
		UploadDate:     formDate(sketch.UploadDate),
		Number:         safeDeref(sketch.Number),
		Popularity:     safeDeref(sketch.Popularity),
		CreatorID:      creatorID,
		CreatorInput:   creatorName,
		EpisodeID:      episodeID,
		EpisodeInput:   episodeName,
		EpisodeStart:   safeDeref(sketch.EpisodeStart),
		SeriesID:       seriesID,
		SeriesInput:    seriesName,
		SeriesPart:     safeDeref(sketch.SeriesPart),
		RecurringID:    recurringID,
		RecurringInput: recurringName,
	}
}

func convertEpisodeToForm(episode *models.Episode) episodeForm {
	var seasonId int
	if episode.Season != nil {
		seasonId = safeDeref(episode.Season.ID)
	}
	return episodeForm{
		ID:            safeDeref(episode.ID),
		Number:        safeDeref(episode.Number),
		Title:         safeDeref(episode.Title),
		URL:           safeDeref(episode.URL),
		AirDate:       formDate(episode.AirDate),
		ThumbnailName: safeDeref(episode.Thumbnail),
		SeasonId:      seasonId,
	}
}

func (app *application) convertFormtoEpisode(form *episodeForm) models.Episode {
	airDate, err := time.Parse(time.DateOnly, form.AirDate)
	episodeAirDate := &airDate
	if err != nil {
		episodeAirDate = nil
	}

	return models.Episode{
		ID:      &form.ID,
		Title:   &form.Title,
		Number:  &form.Number,
		URL:     &form.URL,
		AirDate: episodeAirDate,
		Season: &models.Season{
			ID: &form.SeasonId,
		},
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
		Alias:       &form.Alias,
		BirthDate:   &bdate,
		Professions: &form.Professions,
	}
}

func convertPersontoForm(person *models.Person) personForm {
	return personForm{
		ID:          safeDeref(person.ID),
		First:       safeDeref(person.First),
		Last:        safeDeref(person.Last),
		Alias:       safeDeref(person.Alias),
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
		Aliases:   &form.Aliases,
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
		Aliases:     safeDeref(character.Aliases),
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
		Alias:           &form.Alias,
		URL:             &form.URL,
		EstablishedDate: &date,
	}
}

func convertCreatortoForm(creator *models.Creator) creatorForm {
	return creatorForm{
		ID:              safeDeref(creator.ID),
		Name:            safeDeref(creator.Name),
		Alias:           safeDeref(creator.Alias),
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
		Type:     &form.Type,
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
		Type:          safeDeref(tag.Type),
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
		ID:      &form.ID,
		Name:    &form.Name,
		Aliases: &form.Aliases,
		Slug:    &form.Slug,
	}
}

func (app *application) convertShowtoForm(show *models.Show) showForm {
	return showForm{
		ID:      safeDeref(show.ID),
		Aliases: safeDeref(show.Aliases),
		Name:    safeDeref(show.Name),
		Slug:    safeDeref(show.Slug),
	}
}

func (app *application) convertSeriestoForm(series *models.Series) seriesForm {
	return seriesForm{
		ID:            safeDeref(series.ID),
		Title:         safeDeref(series.Title),
		Description:   safeDeref(series.Description),
		ThumbnailName: safeDeref(series.ThumbnailName),
	}
}

func (app *application) convertFormtoSeries(form *seriesForm) models.Series {
	return models.Series{
		ID:            &form.ID,
		Title:         &form.Title,
		Description:   &form.Description,
		ThumbnailName: &form.ThumbnailName,
	}

}

func (app *application) convertRecurringtoForm(recurring *models.Recurring) recurringForm {
	return recurringForm{
		ID:            safeDeref(recurring.ID),
		Title:         safeDeref(recurring.Title),
		ThumbnailName: safeDeref(recurring.ThumbnailName),
		Description:   safeDeref(recurring.Description),
	}
}

func (app *application) convertFormtoRecurring(form *recurringForm) models.Recurring {
	return models.Recurring{
		ID:            &form.ID,
		Title:         &form.Title,
		Description:   &form.Description,
		ThumbnailName: &form.ThumbnailName,
	}

}

func (app *application) convertFormtoMoment(form *momentForm) models.Moment {
	intTime, _ := models.ParseTimestamp(form.Timestamp)
	return models.Moment{
		ID:        &form.ID,
		Timestamp: &intTime,
		Sketch:    &models.Sketch{ID: &form.SketchID},
	}

}

func (app *application) convertMomenttoForm(moment *models.Moment) momentForm {
	var sketchId int
	if moment.Sketch != nil {
		sketchId = safeDeref(moment.Sketch.ID)
	}

	timestampString := models.SecondsToMMSS(safeDeref(moment.Timestamp))
	return momentForm{
		ID:        safeDeref(moment.ID),
		SketchID:  sketchId,
		Timestamp: timestampString,
	}

}

func (app *application) convertFormtoQuotes(f *quoteForm) []*models.Quote {
	n := len(f.QuoteID)
	quotes := make([]*models.Quote, 0, n)
	for i := range n {
		quotes = append(quotes, &models.Quote{
			ID:         &f.QuoteID[i],
			CastMember: &models.CastMember{ID: &f.CastMemberID[i]},
			Text:       &f.LineText[i],
			Type:       &f.LineType[i],
			Funny:      &f.Funny[i],
			Position:   ptr(i),
		})
	}

	return quotes
}

func (app *application) convertQuotestoForm(sketchId int, momentId int, quotes []*models.Quote) quoteForm {
	f := quoteForm{SketchID: sketchId, MomentID: momentId}
	for _, q := range quotes {
		if q.CastMember != nil {
			var img, text string
			castId := safeDeref(q.CastMember.ID)
			cm, _ := app.cast.GetById(castId)
			if cm == nil {
				castId = 0
			} else {
				img = fmt.Sprintf(
					"%s/cast/profile/small/%s",
					app.baseImgUrl,
					safeDeref(q.CastMember.ProfileImg))

				text = views.PrintCastBlurb(cm)
			}

			f.CastMemberID = append(f.CastMemberID, castId)
			f.CastImageUrl = append(f.CastImageUrl, img)
			f.CastMemberName = append(f.CastMemberName, text)
		}

		f.QuoteID = append(f.QuoteID, safeDeref(q.ID))
		f.LineText = append(f.LineText, safeDeref(q.Text))
		f.LineType = append(f.LineType, views.UppercaseFirst(safeDeref(q.Type)))
		f.Funny = append(f.Funny, views.UppercaseFirst(safeDeref(q.Funny)))
		f.TagCount = append(f.TagCount, len(q.Tags))
	}

	return f
}
