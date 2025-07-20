package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type ShowPage struct {
	ID             int
	ShowName       string
	Image          string
	UpdateUrl      string
	SeasonUrl      string
	SeasonCount    int
	EpisodeCount   int
	SketchCount    int
	SeasonSection  SeasonSelectGallery
	PopularSection *SketchGallery
	CastSection    *PersonGallery
}

func ShowPageView(show *models.Show, popular []*models.Sketch, cast []*models.Person, baseImgUrl string) (*ShowPage, error) {
	page := ShowPage{}
	if show.ID == nil || show.Slug == nil {
		return nil, fmt.Errorf("Show ID and Slug not defined")
	}

	page.ID = *show.ID

	page.ShowName = "Missing Show Name"
	if show.Name != nil {
		page.ShowName = *show.Name
	}

	page.Image = "/static/img/missing-thumbnail.jpg"
	if show.ProfileImg != nil {
		page.Image = fmt.Sprintf("%s/show/%s", baseImgUrl, *show.ProfileImg)
	}

	page.UpdateUrl = fmt.Sprintf("/show/%d/update", *show.ID)
	page.SeasonUrl = fmt.Sprintf("/show/%d/%s/season", *show.ID, *show.Slug)

	page.SeasonCount = len(show.Seasons)
	page.EpisodeCount = getShowEpisodeCount(show)
	page.SketchCount = getShowSketchCount(show)

	page.SeasonSection = SeasonSelectGalleryView(show.Seasons, show.Seasons[0], baseImgUrl, "sub")

	var err error
	page.PopularSection, err = SketchGalleryView(popular, baseImgUrl, baseImgUrl, "sub", 8)
	if err != nil {
		return nil, err
	}

	page.CastSection, err = PersonGalleryView(cast, baseImgUrl)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

type SeasonPage struct {
	ShowName            string
	SeasonNumber        int
	SeasonSelectGallery SeasonSelectGallery
}

func SeasonPageView(show *models.Show, season *models.Season, baseImgUrl string) SeasonPage {
	page := SeasonPage{}
	page.ShowName = "Missing Show Name"
	if show.Name != nil {
		page.ShowName = *show.Name
	}

	if season.Number != nil {
		page.SeasonNumber = *season.Number
	}

	page.SeasonSelectGallery = SeasonSelectGalleryView(show.Seasons, season, baseImgUrl, "Full")
	return page
}

type SeasonSelectGallery struct {
	SeasonUrl      string
	SelectedSeason int
	Seasons        []int
	EpisodeCount   int
	EpisodeGallery EpisodeGallery
}

func SeasonSelectGalleryView(seasons []*models.Season, selected *models.Season, baseImgurl, sectionType string) SeasonSelectGallery {
	gallery := SeasonSelectGallery{}
	if selected.Number != nil {
		gallery.SelectedSeason = *selected.Number
	}

	if selected.ShowId != nil && selected.ShowName != nil {
		gallery.SeasonUrl = fmt.Sprintf(
			"/show/%d/%s/season",
			*selected.ShowId,
			*selected.ShowSlug,
		)
	}
	for _, s := range seasons {
		if s.Number != nil {
			gallery.Seasons = append(gallery.Seasons, *s.Number)
		}
	}

	gallery.EpisodeCount = len(selected.Episodes)
	gallery.EpisodeGallery = EpisodeGalleryView(selected.Episodes, baseImgurl, sectionType)
	return gallery
}

type EpisodePage struct {
	ID               int
	EpisodeTitle     string
	EpisodeInfo      string
	Image            string
	AirDate          string
	ShowName         string
	ShowUrl          string
	ShowImage        string
	SketchCount      int
	Sketches         *SketchGallery
	WatchUrl         string
	WatchImage       string
	UpdateEpisodeUrl string
}

func EpisodePageView(show *models.Show, episode *models.Episode, baseImgUrl string) (*EpisodePage, error) {
	if episode.ID == nil {
		return nil, fmt.Errorf("Episode ID not defined")

	}

	page := EpisodePage{}
	page.ID = *episode.ID

	page.EpisodeTitle = createEpisodeTitle(episode)
	page.EpisodeInfo = seasonEpisodeInfo(episode)
	page.WatchUrl, page.WatchImage = determineEpisodeWatchURL(episode)

	page.Image = "/static/img/missing-thumbnail.jpg"
	if episode.Thumbnail != nil {
		page.Image = fmt.Sprintf(
			"%s/episode/large/%s",
			baseImgUrl,
			*episode.Thumbnail,
		)
	}

	page.AirDate = humanDate(episode.AirDate)
	page.ShowName = "Missing Show"
	if show.Name != nil {
		page.ShowName = *show.Name
	}

	if show.ID != nil && show.Slug != nil {
		page.ShowUrl = fmt.Sprintf(
			"/show/%d/%s",
			*show.ID,
			*show.Slug,
		)
		page.UpdateEpisodeUrl = fmt.Sprintf(
			"/show/%d/update",
			*show.ID,
		)
	}

	page.ShowImage = "/static/img/missing-profile.jpg"
	if show.ProfileImg != nil {
		page.ShowImage = fmt.Sprintf(
			"%s/show/%s",
			baseImgUrl,
			*show.ProfileImg,
		)
	}

	var err error
	page.SketchCount = len(episode.Sketches)
	page.Sketches, err = SketchGalleryView(
		episode.Sketches,
		baseImgUrl,
		"base",
		"full",
		1000,
	)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

type EpisodeGallery struct {
	EpisodeThumbnails []*EpisodeThumbnail
	SectionType       string
}

func EpisodeGalleryView(episodes []*models.Episode, baseImgUrl, sectionType string) EpisodeGallery {
	var episodeThumbnails []*EpisodeThumbnail
	for _, e := range episodes {
		thumbnail := EpisodeThumbnailView(e, baseImgUrl)
		episodeThumbnails = append(episodeThumbnails, thumbnail)
	}

	return EpisodeGallery{
		EpisodeThumbnails: episodeThumbnails,
		SectionType:       sectionType,
	}
}

type EpisodeThumbnail struct {
	Title   string
	Url     string
	Image   string
	AirDate string
	Info    string
}

func EpisodeThumbnailView(episode *models.Episode, baseImgUrl string) *EpisodeThumbnail {
	ep := EpisodeThumbnail{}

	ep.Title = createEpisodeTitle(episode)
	ep.Image = "/static/img/missing-thumbnail.jpg"
	if episode.Thumbnail != nil {
		ep.Image = fmt.Sprintf("%s/episode/%s", baseImgUrl, *episode.Thumbnail)
	}

	ep.Url = seasonEpisodeUrl(episode)
	ep.Info = seasonEpisodeInfo(episode)

	if episode.AirDate != nil {
		ep.AirDate = episode.AirDate.UTC().Format("Jan 2, 2006")
	}

	return &ep
}

type ShowGallery struct {
	ShowCards []*ShowCard
}

type ShowCard struct {
	Name  string
	Url   string
	Image string
}

func ShowGalleryView(shows []*models.Show, baseImgUrl string) (*ShowGallery, error) {
	showGallery := ShowGallery{}

	for _, show := range shows {
		showCard, err := ShowCardView(show, baseImgUrl)
		if err != nil {
			return nil, err
		}

		showGallery.ShowCards = append(showGallery.ShowCards, showCard)
	}

	return &showGallery, nil
}

func ShowCardView(show *models.Show, baseImgUrl string) (*ShowCard, error) {
	card := &ShowCard{}

	if show.ID == nil {
		return nil, fmt.Errorf("Show ID not defined")
	}

	if show.Slug == nil {
		return nil, fmt.Errorf("Show slug not defined")
	}

	card.Name = "Missing Show Name"
	if show.Name != nil {
		card.Name = *show.Name
	}

	card.Url = fmt.Sprintf("/show/%d/%s", *show.ID, *show.Slug)
	card.Image = "/static/img/missing-profile.jpg"
	if show.ProfileImg != nil {
		card.Image = fmt.Sprintf("%s/show/%s", baseImgUrl, *show.ProfileImg)
	}

	return card, nil
}

type FormModal struct {
	Title string
	Form  any
}

type SeasonDropdowns struct {
	ShowID          int
	SeasonDropdowns []SeasonDropdown
}

type SeasonDropdown struct {
	SeasonID     int
	SeasonNumber int
	EpisodeCount int
	EpisodeTable EpisodeTable
}

type EpisodeTable struct {
	SeasonId    int
	EpisodeRows []EpisodeRow
}

type EpisodeRow struct {
	ID           int
	Number       int
	Title        string
	AirDate      string
	ThumbnailUrl string
	SketchCount  int
	EpisodeUrl   string
	SeasonId     int
}

func SeasonDropdownsView(show *models.Show, baseImgUrl string) SeasonDropdowns {
	var dropdowns []SeasonDropdown
	showUrl := fmt.Sprintf("/show/%d/%s", safeDeref(show.ID), safeDeref(show.Slug))
	for _, season := range show.Seasons {
		d := SeasonDropdownView(season, baseImgUrl, showUrl)
		dropdowns = append(dropdowns, d)
	}

	return SeasonDropdowns{
		ShowID:          safeDeref(show.ID),
		SeasonDropdowns: dropdowns,
	}
}

func SeasonDropdownView(season *models.Season, baseImgUrl, showUrl string) SeasonDropdown {
	var d SeasonDropdown
	d.SeasonNumber = safeDeref(season.Number)
	d.EpisodeCount = len(season.Episodes)
	d.EpisodeTable = EpisodeTableView(season, baseImgUrl, showUrl)
	d.SeasonID = safeDeref(season.ID)
	return d
}

func EpisodeTableView(season *models.Season, baseImgUrl, showUrl string) EpisodeTable {
	var rows []EpisodeRow
	for _, episode := range season.Episodes {
		er := EpisodeRow{}
		er.ID = safeDeref(episode.ID)
		er.Number = safeDeref(episode.Number)
		er.Title = safeDeref(episode.Title)
		er.AirDate = humanDate(episode.AirDate)

		er.ThumbnailUrl = fmt.Sprintf("%s/episode/%s", baseImgUrl, safeDeref(episode.Thumbnail))
		er.SketchCount = len(episode.Sketches)
		er.EpisodeUrl = fmt.Sprintf(
			"%s/season/%d/episode/%d",
			showUrl,
			safeDeref(season.Number),
			er.Number,
		)

		er.SeasonId = safeDeref(season.ID)

		rows = append(rows, er)
	}

	return EpisodeTable{
		SeasonId:    safeDeref(season.ID),
		EpisodeRows: rows,
	}
}
