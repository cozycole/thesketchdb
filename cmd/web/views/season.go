package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type SeasonPage struct {
	ShowName            string
	ShowImage           string
	ShowUrl             string
	SeasonNumber        int
	SeasonSelectGallery SeasonSelectGallery
}

func SeasonPageView(show *models.Show, season *models.Season, baseImgUrl string) SeasonPage {
	page := SeasonPage{}
	page.ShowName = "Missing Show Name"
	if show.Name != nil {
		page.ShowName = *show.Name
	}
	page.ShowImage = fmt.Sprintf("%s/show/small/%s", baseImgUrl, safeDeref(show.ProfileImg))
	page.ShowUrl = fmt.Sprintf("/show/%d/%s", safeDeref(show.ID), safeDeref(show.Slug))

	if season.Number != nil {
		page.SeasonNumber = *season.Number
	}

	page.SeasonSelectGallery = SeasonSelectGalleryView(show.Seasons, season, baseImgUrl, "Full")
	return page
}

type SeasonSelectGallery struct {
	SeasonUrl      string
	SelectedSeason int
	Seasons        []SeasonData
	EpisodeCount   int
	EpisodeGallery EpisodeGallery
}

type SeasonData struct {
	Url    string
	Number int
}

func SeasonSelectGalleryView(seasons []*models.Season, selected *models.Season, baseImgurl, sectionType string) SeasonSelectGallery {
	gallery := SeasonSelectGallery{}
	if selected.Number != nil {
		gallery.SelectedSeason = *selected.Number
	}

	if selected.Show != nil && selected.Show.ID != nil {
		gallery.SeasonUrl = fmt.Sprintf(
			"/show/%d/%s/season",
			safeDeref(selected.Show.ID),
			safeDeref(selected.Show.Slug),
		)
	}
	for _, s := range seasons {
		url := fmt.Sprintf(
			"/season/%d/%s",
			safeDeref(s.ID),
			safeDeref(s.Slug),
		)
		if sectionType == "sub" {
			url += "?format=sub"
		}
		seasonData := SeasonData{
			Url:    url,
			Number: safeDeref(s.Number),
		}
		gallery.Seasons = append(gallery.Seasons, seasonData)
	}

	gallery.EpisodeCount = len(selected.Episodes)
	gallery.EpisodeGallery = EpisodeGalleryView(selected.Episodes, baseImgurl, sectionType, false)
	return gallery
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
