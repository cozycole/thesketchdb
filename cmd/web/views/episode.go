package views

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

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

func EpisodePageView(episode *models.Episode, baseImgUrl string) (*EpisodePage, error) {
	if episode.ID == nil {
		return nil, errors.New("Episode ID not defined")
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

	if episode.GetShow() != nil && episode.GetShow().ID != nil {
		page.ShowName = safeDeref(episode.GetShow().Name)
		if page.ShowName == "" {
			page.ShowName = "Missing Show"
		}

		page.UpdateEpisodeUrl = fmt.Sprintf(
			"/show/%d/update",
			safeDeref(episode.GetShow().ID),
		)

		page.ShowUrl = fmt.Sprintf(
			"/show/%d/%s",
			safeDeref(episode.GetShow().ID),
			safeDeref(episode.GetShow().Slug),
		)

		page.ShowImage = "/static/img/missing-profile.jpg"
		page.ShowImage = fmt.Sprintf(
			"%s/show/small/%s",
			baseImgUrl,
			safeDeref(episode.GetShow().ProfileImg),
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
	CountLabel        string
	SectionType       string
}

func EpisodeGalleryView(episodes []*models.EpisodeRef, baseImgUrl, sectionType string, countLabel bool) EpisodeGallery {
	var episodeThumbnails []*EpisodeThumbnail
	for _, e := range episodes {
		thumbnail := EpisodeThumbnailView(e, baseImgUrl)
		episodeThumbnails = append(episodeThumbnails, thumbnail)
	}

	var label string
	if countLabel {
		label = episodeCountLabel(len(episodes))
	}

	return EpisodeGallery{
		EpisodeThumbnails: episodeThumbnails,
		CountLabel:        label,
		SectionType:       sectionType,
	}
}

type EpisodeThumbnail struct {
	Title       string
	Url         string
	Image       string
	LargeImage  string
	MediumImage string
	SmallImage  string
	AirDate     string
	Info        string
}

func EpisodeThumbnailView(episode *models.EpisodeRef, baseImgUrl string) *EpisodeThumbnail {
	ep := EpisodeThumbnail{}

	ep.Title = createEpisodeTitle(episode)
	ep.Image = "/static/img/missing-thumbnail.jpg"
	ep.LargeImage = "/static/img/missing-thumbnail.jpg"
	ep.MediumImage = "/static/img/missing-thumbnail.jpg"
	ep.SmallImage = "/static/img/missing-thumbnail.jpg"
	if episode.Thumbnail != nil {
		ep.Image = fmt.Sprintf("%s/episode/small/%s", baseImgUrl, *episode.Thumbnail)
		ep.LargeImage = fmt.Sprintf("%s/episode/large/%s", baseImgUrl, *episode.Thumbnail)
		ep.MediumImage = fmt.Sprintf("%s/episode/medium/%s", baseImgUrl, *episode.Thumbnail)
		ep.SmallImage = fmt.Sprintf("%s/episode/small/%s", baseImgUrl, *episode.Thumbnail)
	}

	ep.Url = fmt.Sprintf(
		"/episode/%d/%s",
		safeDeref(episode.ID),
		safeDeref(episode.Slug),
	)

	ep.Info = seasonEpisodeInfo(episode)

	if episode.AirDate != nil {
		ep.AirDate = episode.AirDate.UTC().Format("Jan 2, 2006")
	}

	return &ep
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

func EpisodeTableView(season *models.Season, baseImgUrl, showUrl string) EpisodeTable {
	var rows []EpisodeRow
	for _, episode := range season.Episodes {
		er := EpisodeRow{}
		er.ID = safeDeref(episode.ID)
		er.Number = safeDeref(episode.Number)
		er.Title = safeDeref(episode.Title)
		er.AirDate = humanDate(episode.AirDate)

		er.ThumbnailUrl = fmt.Sprintf("%s/episode/small/%s", baseImgUrl, safeDeref(episode.Thumbnail))
		er.SketchCount = safeDeref(episode.SketchCount)
		er.EpisodeUrl = fmt.Sprintf(
			"/episode/%d/%s",
			safeDeref(episode.ID),
			safeDeref(episode.Slug),
		)

		er.SeasonId = safeDeref(season.ID)

		rows = append(rows, er)
	}

	return EpisodeTable{
		SeasonId:    safeDeref(season.ID),
		EpisodeRows: rows,
	}
}
