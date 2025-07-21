package views

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

const YOUTUBE_URL = "https://www.youtube.com/watch?v=%s&t=%ds"

type SketchPage struct {
	SketchID      int
	Title         string
	Description   string
	Image         string
	UpdateUrl     string
	YoutubeId     string
	YoutubeUrl    string
	Date          string
	Liked         bool
	CreatorName   string
	CreatorImage  string
	CreatorUrl    string
	SeasonNumber  int
	SeasonUrl     string
	EpisodeNumber int
	EpisodeUrl    string
	SketchNumber  int
	StartTime     int
	SeriesPart    int
	SeriesTitle   string
	SeriesUrl     string
	InSeries      bool
	Cast          CastGallery
	Tags          []*Tag
}

func SketchPageView(sketch *models.Sketch, tags []*models.Tag, baseImgUrl string) (*SketchPage, error) {
	page := SketchPage{}
	if sketch.ID == nil {
		return nil, fmt.Errorf("Sketch ID not defined")
	}

	if sketch.Slug == nil {
		return nil, fmt.Errorf("Sketch slug not defined")
	}

	page.SketchID = *sketch.ID

	page.Image = "/static/img/missing-thumbnail.jpg"
	if sketch.ThumbnailName != nil {
		page.Image = fmt.Sprintf("%s/sketch/large/%s", baseImgUrl, *sketch.ThumbnailName)
	}

	page.Title = safeDeref(sketch.Title)
	if page.Title == "" {
		page.Title = "Missing Title"
	}

	page.Description = safeDeref(sketch.Description)

	if sketch.UploadDate != nil {
		page.Date = sketch.UploadDate.UTC().Format("Jan 2, 2006")
	}

	if sketch.YoutubeID != nil {
		page.YoutubeId = *sketch.YoutubeID
		page.YoutubeUrl = fmt.Sprintf(YOUTUBE_URL, *sketch.YoutubeID, 0)
	} else if sketch.Episode != nil {
		page.YoutubeId = safeDeref(sketch.Episode.YoutubeID)
		page.YoutubeUrl = fmt.Sprintf(YOUTUBE_URL, page.YoutubeId, safeDeref(sketch.EpisodeStart))
		page.StartTime = safeDeref(sketch.EpisodeStart)
	}

	if sketch.Liked != nil {
		page.Liked = *sketch.Liked
	}

	page.UpdateUrl = fmt.Sprintf("/sketch/%d/update", *sketch.ID)
	page.Cast, _ = CastGalleryView(sketch.Cast, baseImgUrl)
	page.Tags = TagsView(tags)

	if sketch.Show != nil && sketch.Show.ID != nil {
		if sketch.Show.Name != nil {
			page.CreatorName = *sketch.Show.Name
		}

		if sketch.Show.ID != nil && sketch.Show.Slug != nil {
			page.CreatorUrl = fmt.Sprintf("/show/%d/%s", *sketch.Show.ID, *sketch.Show.Slug)
		}

		if sketch.Show.ProfileImg != nil {
			page.CreatorImage = fmt.Sprintf("%s/show/%s", baseImgUrl, *sketch.Show.ProfileImg)
		} else {
			page.CreatorImage = fmt.Sprintf("%s/missing-profile.jpg", baseImgUrl)
		}

		if sketch.Season != nil {
			page.SeasonNumber = safeDeref(sketch.Season.Number)
			page.SeasonUrl = fmt.Sprintf(
				"/season/%d/%s",
				safeDeref(sketch.Season.ID),
				safeDeref(sketch.Season.Slug),
			)
		}

		if sketch.Episode != nil && safeDeref(sketch.Episode.ID) != 0 {
			ep := sketch.Episode
			page.EpisodeNumber = safeDeref(ep.Number)
			page.EpisodeUrl = fmt.Sprintf(
				"/episode/%d/%s",
				safeDeref(ep.ID),
				safeDeref(ep.Slug),
			)
		}

		page.SketchNumber = safeDeref(sketch.Number)

	} else if sketch.Creator != nil && sketch.Creator.ID != nil {
		if sketch.Creator.Name != nil {
			page.CreatorName = *sketch.Creator.Name
		}

		if sketch.Creator.ID != nil && sketch.Creator.Slug != nil {
			page.CreatorUrl = fmt.Sprintf("/creator/%d/%s", *sketch.Creator.ID, *sketch.Creator.Slug)
		}

		if sketch.Creator.ProfileImage != nil {
			page.CreatorImage = fmt.Sprintf("%s/creator/%s", baseImgUrl, *sketch.Creator.ProfileImage)
		}
	} else {
		page.CreatorName = "Missing Creator"
		page.CreatorImage = fmt.Sprintf("%s/missing-profile.jpg", baseImgUrl)
	}

	if sketch.Series != nil {
		page.SeriesTitle = safeDeref(sketch.Series.Title)
		page.SeriesUrl = fmt.Sprintf(
			"/series/%d/%s",
			safeDeref(sketch.Series.ID),
			safeDeref(sketch.Series.Slug),
		)
		page.SeriesPart = safeDeref(sketch.SeriesPart)

		page.InSeries = page.SeriesTitle != "" &&
			sketch.Series.ID != nil &&
			safeDeref(sketch.SeriesPart) != 0
	}

	return &page, nil
}

type SketchGallery struct {
	Sketches []*SketchThumbnail
	// if it's in a subsection, grid col breakpoints need to be changed
	SectionType string
	SeeMore     bool
	SeeMoreUrl  string
}

type SketchThumbnail struct {
	Title        string
	Url          string
	YoutubeUrl   string
	Image        string
	LargeImage   string
	Date         string
	Liked        bool
	CreatorName  string
	CreatorImage string
	CreatorUrl   string
	CreatorInfo  string
}

func SketchGalleryView(
	sketches []*models.Sketch,
	baseImgUrl,
	thumbnailType,
	sectionType string,
	maxResults int,
) (*SketchGallery, error) {
	sketchViews, err := SketchThumbnailsView(sketches, baseImgUrl, thumbnailType)
	if err != nil {
		return nil, err
	}

	return &SketchGallery{
		Sketches:    sketchViews,
		SectionType: sectionType,
	}, nil
}

func FeaturedSketchesView(sketches []*models.Sketch, baseImgUrl string) ([]*SketchThumbnail, error) {
	var sketchViews []*SketchThumbnail
	for _, sketch := range sketches {
		sketchView, err := SketchThumbnailView(sketch, baseImgUrl, "")
		if err != nil {
			return nil, err
		}

		sketchView.CreatorInfo = printCast(sketch.Cast)

		sketchViews = append(sketchViews, sketchView)
	}

	return sketchViews, nil
}

func SketchThumbnailsView(sketches []*models.Sketch, baseImgUrl string, thumbnailType string) ([]*SketchThumbnail, error) {
	var sketchViews []*SketchThumbnail
	for _, sketch := range sketches {
		sketchView, err := SketchThumbnailView(sketch, baseImgUrl, thumbnailType)
		if err != nil {
			return nil, err
		}

		sketchViews = append(sketchViews, sketchView)
	}
	return sketchViews, nil
}

func SketchThumbnailView(sketch *models.Sketch, baseImgUrl string, thumbnailType string) (*SketchThumbnail, error) {
	sketchView := &SketchThumbnail{}
	if sketch.ID == nil {
		return nil, fmt.Errorf("Sketch ID not defined")
	}

	if sketch.Slug == nil {
		return nil, fmt.Errorf("Sketch slug not defined")
	}

	if sketch.Title != nil && *sketch.Title != "" {
		sketchView.Title = *sketch.Title
	} else {
		sketchView.Title = "Untitled Sketch"
	}

	sketchView.Url = fmt.Sprintf("/sketch/%d/%s", *sketch.ID, *sketch.Slug)

	if sketch.YoutubeID != nil && len(*sketch.YoutubeID) == 11 {
		sketchView.YoutubeUrl = fmt.Sprintf("www.youtube.com/watch?v=%s", *sketch.YoutubeID)
	}

	if strings.ToUpper(thumbnailType) == "CAST" && safeDeref(sketch.CastThumbnail) != "" {
		sketchView.Image = fmt.Sprintf("%s/cast/thumbnail/%s", baseImgUrl, safeDeref(sketch.CastThumbnail))
		// there are no Large cast images
		sketchView.LargeImage = fmt.Sprintf("%s/sketch/large/%s", baseImgUrl, safeDeref(sketch.ThumbnailName))
	} else if sketch.ThumbnailName != nil {
		sketchView.Image = fmt.Sprintf("%s/sketch/%s", baseImgUrl, *sketch.ThumbnailName)
		sketchView.LargeImage = fmt.Sprintf("%s/sketch/large/%s", baseImgUrl, *sketch.ThumbnailName)
	} else {
		sketchView.Image = fmt.Sprintf("%s/missing-thumbnail.jpg", baseImgUrl)
		sketchView.LargeImage = fmt.Sprintf("%s/missing-thumbnail.jpg", baseImgUrl)
	}

	if sketch.UploadDate != nil {
		sketchView.Date = sketch.UploadDate.UTC().Format("Jan 2, 2006")
	}

	if sketch.Show != nil && sketch.Show.ID != nil {
		if sketch.Show.Name != nil {
			sketchView.CreatorName = *sketch.Show.Name
		}

		if sketch.Show.ID != nil && sketch.Show.Slug != nil {
			sketchView.CreatorUrl = fmt.Sprintf("/show/%d/%s", *sketch.Show.ID, *sketch.Show.Slug)
		}

		if sketch.Show.ProfileImg != nil {
			sketchView.CreatorImage = fmt.Sprintf("%s/show/%s", baseImgUrl, *sketch.Show.ProfileImg)
		} else {
			sketchView.CreatorImage = fmt.Sprintf("%s/missing-profile.jpg", baseImgUrl)
		}

		var season, episode, number int
		if sketch.Season != nil && sketch.Season.Number != nil {
			season = *sketch.Season.Number
		}

		if sketch.Episode != nil && sketch.Episode.Number != nil {
			episode = *sketch.Episode.Number
		}

		if sketch.Number != nil {
			number = *sketch.Number
		}

		sketchView.CreatorInfo = fmt.Sprintf("S%d · E%d · #%d", season, episode, number)
	} else if sketch.Creator != nil && sketch.Creator.ID != nil {
		if sketch.Creator.Name != nil {
			sketchView.CreatorName = *sketch.Creator.Name
		}

		if sketch.Creator.ID != nil && sketch.Creator.Slug != nil {
			sketchView.CreatorUrl = fmt.Sprintf("/creator/%d/%s", *sketch.Creator.ID, *sketch.Creator.Slug)
		}

		if sketch.Creator.ProfileImage != nil {
			sketchView.CreatorImage = fmt.Sprintf("%s/creator/%s", baseImgUrl, *sketch.Creator.ProfileImage)
		}
	} else {
		sketchView.CreatorName = "Missing Creator"
		sketchView.CreatorImage = fmt.Sprintf("%s/missing-profile.jpg", baseImgUrl)
	}

	return sketchView, nil
}

type SketchCatalog struct {
	ResultCountLabel string
	CatalogFilter    SketchCatalogFilter
	CatalogResult    SketchCatalogResult
}

func SketchCatalogView(
	results *models.SearchResult,
	currentPage int,
	totalPages int,
	htmxRequest bool,
	baseImgUrl string,
) (*SketchCatalog, error) {
	sketchCatalogResult, err := SketchCatalogResultView(
		results,
		currentPage,
		totalPages,
		htmxRequest,
		baseImgUrl,
	)
	if err != nil {
		return nil, err
	}

	sketchCatalogFilter, err := SketchCatalogFilterView(
		results.Filter,
		baseImgUrl,
	)
	if err != nil {
		return nil, err
	}

	return &SketchCatalog{
		ResultCountLabel: sketchCountLabel(results.TotalSketchCount),
		CatalogFilter:    *sketchCatalogFilter,
		CatalogResult:    *sketchCatalogResult,
	}, nil
}

type SketchCatalogResult struct {
	HasResults           bool
	ResultCountLabel     string
	IsHtmxRequest        bool
	SketchResultsGallery SketchGallery
	Pages                []*PaginationItem
}

func SketchCatalogResultView(
	results *models.SearchResult,
	currentPage int,
	totalPages int,
	htmxRequest bool,
	baseImgUrl string,
) (*SketchCatalogResult, error) {
	thumbnailType := "Base"
	if len(results.Filter.People) == 1 || len(results.Filter.Characters) == 1 {
		thumbnailType = "Cast"
	}

	sketches, err := SketchThumbnailsView(
		results.SketchResults,
		baseImgUrl,
		thumbnailType,
	)
	if err != nil {
		return nil, err
	}

	pages, err := buildPagination(
		currentPage,
		totalPages,
		"/catalog/sketches",
		results.Filter,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("current page", currentPage, "totalPages", totalPages, results.Filter.ParamsString())
	for _, p := range pages {
		fmt.Printf("PAGES: %+v\n", p)
	}

	labelString := "%d Sketch"
	if results.TotalSketchCount != 1 {
		labelString += "es"
	}

	labelString = fmt.Sprintf(labelString, results.TotalSketchCount)

	return &SketchCatalogResult{
		HasResults:           len(results.SketchResults) != 0,
		IsHtmxRequest:        htmxRequest,
		ResultCountLabel:     labelString,
		SketchResultsGallery: SketchGallery{Sketches: sketches, SectionType: "full"},
		Pages:                pages,
	}, nil
}

type SketchCatalogFilter struct {
	SortOptions            []SortOption
	SelectedPeopleJSON     string
	SelectedCreatorsJSON   string
	SelectedShowsJSON      string
	SelectedCharactersJSON string
	SelectedTagsJSON       string
}

type SortOption struct {
	Value    string
	Label    string
	Selected bool
}

func SketchCatalogFilterView(filter *models.Filter, baseUrl string) (*SketchCatalogFilter, error) {
	var view SketchCatalogFilter
	sortBy := filter.SortBy
	view.SortOptions = []SortOption{
		{Value: "popular", Label: "Popular", Selected: sortBy == "popular"},
		{Value: "latest", Label: "Latest", Selected: sortBy == "latest"},
		{Value: "oldest", Label: "Oldest", Selected: sortBy == "oldest"},
		{Value: "az", Label: "A-Z", Selected: sortBy == "az"},
		{Value: "za", Label: "Z-A", Selected: sortBy == "za"},
	}
	var err error
	if view.SelectedPeopleJSON, err = PeopleSelectedJSON(filter.People, baseUrl); err != nil {
		return nil, err
	}

	if view.SelectedCreatorsJSON, err = CreatorsSelectedJSON(filter.Creators, baseUrl); err != nil {
		return nil, err
	}

	if view.SelectedCharactersJSON, err = CharactersSelectedJSON(filter.Characters, baseUrl); err != nil {
		return nil, err
	}
	if view.SelectedShowsJSON, err = ShowsSelectedJSON(filter.Shows, baseUrl); err != nil {
		return nil, err
	}
	if view.SelectedTagsJSON, err = TagsSelectedJSON(filter.Tags); err != nil {
		return nil, err
	}

	return &view, nil
}

type SelectedItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

func PeopleSelectedJSON(people []*models.Person, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(people))
	for _, p := range people {

		var image string
		if p.ProfileImg != nil {
			image = fmt.Sprintf("%s/person/%s", baseURL, *p.ProfileImg)
		}

		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*p.ID),
			Name:  PrintPersonName(p),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func CreatorsSelectedJSON(creators []*models.Creator, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(creators))
	for _, c := range creators {

		var image string
		if c.ProfileImage != nil {
			image = fmt.Sprintf("%s/creator/%s", baseURL, *c.ProfileImage)
		}
		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*c.ID),
			Name:  safeDeref(c.Name),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func ShowsSelectedJSON(shows []*models.Show, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(shows))
	for _, s := range shows {

		var image string
		if s.ProfileImg != nil {
			image = fmt.Sprintf("%s/show/%s", baseURL, *s.ProfileImg)
		}
		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*s.ID),
			Name:  safeDeref(s.Name),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func CharactersSelectedJSON(characters []*models.Character, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(characters))
	for _, c := range characters {

		var image string
		if c.Image != nil {
			image = fmt.Sprintf("%s/character/%s", baseURL, *c.Image)
		}
		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*c.ID),
			Name:  safeDeref(c.Name),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func TagsSelectedJSON(tags []*models.Tag) (string, error) {
	items := make([]SelectedItem, 0, len(tags))
	for _, t := range tags {

		var name string
		if t.Category != nil && safeDeref(t.Category.Name) != "" {
			name = fmt.Sprintf("%s / ", safeDeref(t.Category.Name))
		}
		name += safeDeref(t.Name)
		items = append(items, SelectedItem{
			ID:   strconv.Itoa(*t.ID),
			Name: name,
		})
	}

	return buildSelectedJSON(items)
}

func buildSelectedJSON(items []SelectedItem) (string, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
