package views

import (
	"errors"
	"fmt"
	"html/template"

	"sketchdb.cozycole.net/internal/domain/sketches"
	"sketchdb.cozycole.net/internal/models"
)

type ShowTab string

const (
	ShowTabHome       ShowTab = "home"
	ShowTabSketches   ShowTab = "sketches"
	ShowTabSeasons    ShowTab = "seasons"
	ShowTabExtras     ShowTab = "extras"
	ShowTabCast       ShowTab = "cast"
	ShowTabCharacters ShowTab = "characters"
	ShowTabQuotes     ShowTab = "quotes"
)

type baseShowLayout struct {
	ID               int
	Slug             string
	ShowName         string
	Image            string
	UpdateUrl        string
	SketchCatalogUrl string
	SeasonCount      int
	EpisodeCount     int
	SketchCount      int

	ActiveTab ShowTab
}

type ShowHomePage struct {
	baseShowLayout
	About          template.HTML
	SeasonsUrl     string
	SeasonSection  SeasonSelectGallery
	PopularSection *SketchGallery
	CastSection    *PersonGallery
}

func baseShowLayoutView(show *models.Show, activeTab ShowTab, baseImgUrl string) (baseShowLayout, error) {
	base := baseShowLayout{}
	if show.ID == nil || show.Slug == nil {
		return base, errors.New("Show ID and Slug not defined")
	}

	base.ID = *show.ID
	base.Slug = *show.Slug

	base.ShowName = "Missing Show Name"
	if show.Name != nil {
		base.ShowName = *show.Name
	}

	base.Image = "/static/img/missing-profile.jpg"
	if show.ProfileImg != nil {
		base.Image = fmt.Sprintf("%s/show/medium/%s", baseImgUrl, *show.ProfileImg)
	}

	base.UpdateUrl = fmt.Sprintf("/show/%d/update", *show.ID)
	base.SketchCatalogUrl = fmt.Sprintf("/catalog/sketches?show=%d", *show.ID)

	base.SeasonCount = len(show.Seasons)
	base.EpisodeCount = getShowEpisodeCount(show)
	base.SketchCount = getShowSketchCount(show)
	base.ActiveTab = activeTab
	return base, nil
}

func ShowHomePageView(show *models.Show, popular []*models.SketchRef, cast []*models.Person, baseImgUrl string) (*ShowHomePage, error) {
	base, err := baseShowLayoutView(show, ShowTabHome, baseImgUrl)
	if err != nil {
		return nil, err
	}

	page := ShowHomePage{
		baseShowLayout: base,
	}

	page.About = template.HTML(safeDeref(show.About))
	page.SeasonSection = SeasonSelectGalleryView(show.Seasons, show.Seasons[0], baseImgUrl, "sub")

	popularPageSize := 12
	page.PopularSection, err = SketchGalleryView(
		popular, baseImgUrl, baseImgUrl, "sub")
	if err != nil {
		return nil, err
	}

	if len(popular) == popularPageSize {
		page.PopularSection.SeeMore = true
		page.PopularSection.SeeMoreUrl = fmt.Sprintf(
			"/show/%d/%s/sketches", page.ID, page.Slug,
		)
	}

	page.CastSection, err = PersonGalleryView(cast, baseImgUrl)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

type ShowSketchesPage struct {
	baseShowLayout
	SketchResultsGallery SketchGallery
	HasResults           bool
	Pages                []*PaginationItem
}

func ShowSketchesPageView(
	show *models.Show,
	results sketches.SketchListResult,
	baseImgUrl string,
) (*ShowSketchesPage, error) {
	base, err := baseShowLayoutView(show, ShowTabSketches, baseImgUrl)
	if err != nil {
		return nil, err
	}

	page := ShowSketchesPage{
		baseShowLayout: base,
	}

	sketches, err := SketchThumbnailsView(
		results.Sketches,
		baseImgUrl,
		"Base",
		false,
	)
	if err != nil {
		return nil, err
	}

	pagination, err := buildPagination(
		results.Metadata.CurrentPage,
		results.Metadata.TotalPages,
		fmt.Sprintf("/show/%d/%show/sketches", page.ID, page.Slug),
		results.Filter,
	)

	page.SketchResultsGallery = SketchGallery{Sketches: sketches, SectionType: "full"}
	page.Pages = pagination
	page.HasResults = len(sketches) > 0
	return &page, nil
}

type ShowSeasonsPage struct {
	baseShowLayout
	SeasonSection SeasonSelectGallery
}

func ShowSeasonsPageView(show *models.Show, baseImgUrl string) (*ShowSeasonsPage, error) {
	base, err := baseShowLayoutView(show, ShowTabSeasons, baseImgUrl)
	if err != nil {
		return nil, err
	}

	page := ShowSeasonsPage{
		baseShowLayout: base,
	}

	var season *models.Season
	if len(show.Seasons) > 0 {
		season = show.Seasons[0]
	}
	page.SeasonSection = SeasonSelectGalleryView(show.Seasons, season, baseImgUrl, "sub")

	return &page, nil
}

type ShowExtrasPage struct {
	baseShowLayout
	Groupings []Grouping
}

func ShowExtrasPageView(show *models.Show, groupings []*models.Grouping, baseImgUrl string) (*ShowExtrasPage, error) {
	base, err := baseShowLayoutView(show, ShowTabExtras, baseImgUrl)
	if err != nil {
		return nil, err
	}

	page := ShowExtrasPage{
		baseShowLayout: base,
	}

	for _, g := range groupings {
		page.Groupings = append(page.Groupings, GroupingView(g, baseImgUrl))
	}

	return &page, nil
}

type ShowCastPage struct {
	baseShowLayout
	CastSection *PersonGallery
}

func ShowCastPageView(show *models.Show, cast []*models.Person, baseImgUrl string) (*ShowHomePage, error) {
	base, err := baseShowLayoutView(show, ShowTabCast, baseImgUrl)
	if err != nil {
		return nil, err
	}

	page := ShowHomePage{
		baseShowLayout: base,
	}
	page.CastSection, err = PersonGalleryView(cast, baseImgUrl)
	if err != nil {
		return nil, err
	}

	page.SeasonSection = SeasonSelectGalleryView(show.Seasons, show.Seasons[0], baseImgUrl, "sub")
	return &page, nil
}

type ShowGallery struct {
	Cards []*Card
}

func ShowGalleryView(shows []*models.Show, baseImgUrl string) (*ShowGallery, error) {
	showGallery := ShowGallery{}

	for _, show := range shows {
		showCard, err := ShowCardView(show, baseImgUrl)
		if err != nil {
			return nil, err
		}

		showGallery.Cards = append(showGallery.Cards, showCard)
	}

	return &showGallery, nil
}

func ShowCardView(show *models.Show, baseImgUrl string) (*Card, error) {
	card := &Card{}

	if show.ID == nil {
		return nil, errors.New("Show ID not defined")
	}

	if show.Slug == nil {
		return nil, errors.New("Show slug not defined")
	}

	card.Title = safeDeref(show.Name)

	card.Url = fmt.Sprintf("/show/%d/%s", *show.ID, *show.Slug)
	card.ImageUrl = "/static/img/missing-profile.jpg"
	if show.ProfileImg != nil {
		card.ImageUrl = fmt.Sprintf("%s/show/medium/%s", baseImgUrl, *show.ProfileImg)
	}

	return card, nil
}

type FormModal struct {
	Title string
	Form  any
}
