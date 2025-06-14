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
	ID           int
	EpisodeTitle string
	EpisodeInfo  string
	Image        string
	AirDate      string
	ShowName     string
	ShowUrl      string
	ShowImage    string
	SketchCount  int
	Sketches     *SketchGallery
}

func EpisodePageView(show *models.Show, episode *models.Episode, baseImgUrl string) (*EpisodePage, error) {
	if episode.ID == nil {
		return nil, fmt.Errorf("Episode ID not defined")

	}

	page := EpisodePage{}
	page.ID = *episode.ID

	page.EpisodeTitle = createEpisodeTitle(episode)
	page.EpisodeInfo = seasonEpisodeInfo(episode)

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
