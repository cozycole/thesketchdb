package views

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/internal/domain/sketches"
	"sketchdb.cozycole.net/internal/models"
)

const YOUTUBE_URL = "https://www.youtube.com/watch?v=%s&t=%ds"

type SketchPage struct {
	SketchID       int
	Title          string
	Description    string
	Image          string
	UpdateUrl      string
	YoutubeId      string
	YoutubeUrl     string
	Date           string
	Liked          bool
	CreatorName    string
	CreatorImage   string
	CreatorUrl     string
	SeasonNumber   int
	SeasonUrl      string
	EpisodeNumber  int
	EpisodeUrl     string
	SketchNumber   int
	StartTime      int
	InSeries       bool
	SeriesPart     int
	SeriesTitle    string
	SeriesUrl      string
	InRecurring    bool
	RecurringTitle string
	RecurringUrl   string
	Rating         SketchRating
	Cast           CastGallery
	Quotes         []Quote
	Tags           []*Tag
}

func SketchPageView(
	sketch *models.Sketch,
	quotes []*models.Quote,
	tags []*models.Tag,
	userSketchInfo *models.UserSketchInfo,
	baseImgUrl string) (*SketchPage, error) {
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

	page.UpdateUrl = fmt.Sprintf("/admin/sketch/%d", *sketch.ID)
	page.Cast, _ = CastGalleryView(sketch.Cast, baseImgUrl)
	page.Tags = TagsView(tags)

	if sketch.Episode != nil && sketch.Episode.ID != nil {
		if sketch.Episode != nil && safeDeref(sketch.Episode.ID) != 0 {
			ep := sketch.Episode
			page.EpisodeNumber = safeDeref(ep.Number)
			page.EpisodeUrl = fmt.Sprintf(
				"/episode/%d/%s",
				safeDeref(ep.ID),
				safeDeref(ep.Slug),
			)
			season := safeDeref(ep.Season)
			if season.ID != nil {
				page.SeasonNumber = safeDeref(season.Number)
				page.SeasonUrl = fmt.Sprintf(
					"/season/%d/%s",
					safeDeref(season.ID),
					safeDeref(season.Slug),
				)

				show := safeDeref(season.Show)
				if show.ID != nil {
					page.CreatorName = safeDeref(show.Name)
					if show.ID != nil && show.Slug != nil {
						page.CreatorUrl = fmt.Sprintf("/show/%d/%s", *show.ID, *show.Slug)
					}
					if show.ProfileImg != nil {
						page.CreatorImage = fmt.Sprintf("%s/show/small/%s", baseImgUrl, *show.ProfileImg)
					} else {
						page.CreatorImage = fmt.Sprintf("%s/missing-profile.jpg", baseImgUrl)
					}
				}
			}
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
			page.CreatorImage = fmt.Sprintf("%s/creator/small/%s", baseImgUrl, *sketch.Creator.ProfileImage)
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

	if sketch.Recurring != nil {
		page.RecurringTitle = safeDeref(sketch.Recurring.Title)
		page.RecurringUrl = fmt.Sprintf(
			"/recurring/%d/%s",
			safeDeref(sketch.Recurring.ID),
			safeDeref(sketch.Recurring.Slug),
		)

		page.InRecurring = page.RecurringTitle != "" &&
			sketch.Recurring.ID != nil
	}

	page.Quotes = SketchQuoteSection(quotes, baseImgUrl)
	page.Rating = SketchRatingView(userSketchInfo, sketch)
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
	SmallImage   string
	MediumImage  string
	LargeImage   string
	Date         string
	Liked        bool
	CreatorName  string
	CreatorImage string
	CreatorUrl   string
	CreatorInfo  string
	Rating       string
	InCarousel   bool
}

func SketchGalleryView(
	sketches []*models.SketchRef,
	baseImgUrl,
	thumbnailType,
	sectionType string,
	maxResults int,
) (*SketchGallery, error) {
	sketchViews, err := SketchThumbnailsView(sketches, baseImgUrl, thumbnailType, false)
	if err != nil {
		return nil, err
	}

	return &SketchGallery{
		Sketches:    sketchViews,
		SectionType: sectionType,
	}, nil
}

// for sketch displays where either it's a carousel or a grid to carousel (carousel on mobile)
func SketchCarouselView(
	sketches []*models.SketchRef,
	baseImgUrl,
	thumbnailType,
	sectionType string,
	maxResults int,
) (*SketchGallery, error) {
	sketchViews, err := SketchThumbnailsView(sketches, baseImgUrl, thumbnailType, true)
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
		sketchRef := convertSketchToSketchRef(sketch)

		sketchView, err := SketchThumbnailView(&sketchRef, baseImgUrl, "", false)
		if err != nil {
			return nil, err
		}

		sketchView.CreatorInfo = printCast(sketch.Cast)

		sketchViews = append(sketchViews, sketchView)
	}

	return sketchViews, nil
}

// this function is only for the Featured Sketches section
func convertSketchToSketchRef(s *models.Sketch) models.SketchRef {
	creator := models.CreatorRef{}
	if s.Creator != nil && s.Creator.ID != nil {
		creator.ID = s.Creator.ID
		creator.Slug = s.Creator.Slug
		creator.Name = s.Creator.Name
		creator.ProfileImage = s.Creator.ProfileImage
	}

	ep := models.EpisodeRef{}
	se := models.SeasonRef{}
	sh := models.ShowRef{}
	if s.Episode != nil && s.Episode.ID != nil {
		ep.ID = s.Episode.ID
		ep.Slug = s.Episode.Slug
		ep.Number = s.Episode.Number
		ep.AirDate = s.Episode.AirDate
		ep.Season = &se

		ogSeason := s.Episode.Season
		if ogSeason != nil && ogSeason.ID != nil {
			se.ID = ogSeason.ID
			se.Slug = ogSeason.Slug
			se.Number = ogSeason.Number
			se.Show = &sh

			ogShow := ogSeason.Show
			if ogShow != nil && ogShow.ID != nil {
				sh.ID = ogShow.ID
				sh.Slug = ogShow.Slug
				sh.Name = ogShow.Name
				sh.ProfileImg = ogShow.ProfileImg
			}
		}
	}
	return models.SketchRef{
		ID:            s.ID,
		Slug:          s.Slug,
		Title:         s.Title,
		Thumbnail:     s.ThumbnailName,
		CastThumbnail: s.CastThumbnail,
		UploadDate:    s.UploadDate,
		Number:        s.Number,
		Rating:        s.Rating,
		Episode:       &ep,
		Creator:       &creator,
	}
}

func SketchThumbnailsView(sketches []*models.SketchRef, baseImgUrl string, thumbnailType string, inCarousel bool) ([]*SketchThumbnail, error) {
	var sketchViews []*SketchThumbnail
	for _, sketch := range sketches {
		sketchView, err := SketchThumbnailView(sketch, baseImgUrl, thumbnailType, inCarousel)
		if err != nil {
			return nil, err
		}

		sketchViews = append(sketchViews, sketchView)
	}
	return sketchViews, nil
}

func SketchThumbnailView(sketch *models.SketchRef, baseImgUrl string, thumbnailType string, inCarousel bool) (*SketchThumbnail, error) {
	if sketch == nil || sketch.ID == nil || sketch.Slug == nil {
		return nil, fmt.Errorf("Sketch ID and slug are not defined")
	}

	sketchView := &SketchThumbnail{}
	if sketch.Title != nil && *sketch.Title != "" {
		sketchView.Title = *sketch.Title
	} else {
		sketchView.Title = "Untitled Sketch"
	}

	sketchView.Url = fmt.Sprintf("/sketch/%d/%s", *sketch.ID, *sketch.Slug)

	sketchView.InCarousel = inCarousel

	if safeDeref(sketch.Thumbnail) != "" && safeDeref(sketch.Thumbnail) != "missing-thumbnail.jpg" {
		sketchView.SmallImage = fmt.Sprintf("%s/sketch/small/%s", baseImgUrl, safeDeref(sketch.Thumbnail))
		sketchView.MediumImage = fmt.Sprintf("%s/sketch/medium/%s", baseImgUrl, safeDeref(sketch.Thumbnail))
		sketchView.LargeImage = fmt.Sprintf("%s/sketch/large/%s", baseImgUrl, safeDeref(sketch.Thumbnail))
		sketchView.Image = sketchView.SmallImage
	} else {
		sketchView.Image = "/static/img/missing-thumbnail.jpg"
		sketchView.SmallImage = "/static/img/missing-thumbnail.jpg"
		sketchView.MediumImage = "/static/img/missing-thumbnail.jpg"
		sketchView.LargeImage = "/static/img/missing-thumbnail.jpg"
	}

	if strings.ToUpper(thumbnailType) == "CAST" && safeDeref(sketch.CastThumbnail) != "" {
		sketchView.SmallImage = fmt.Sprintf("%s/cast/thumbnail/small/%s", baseImgUrl, safeDeref(sketch.CastThumbnail))
		sketchView.MediumImage = fmt.Sprintf("%s/cast/thumbnail/medium/%s", baseImgUrl, safeDeref(sketch.CastThumbnail))
		sketchView.Image = fmt.Sprintf("%s/cast/thumbnail/small/%s", baseImgUrl, safeDeref(sketch.CastThumbnail))
	}

	if sketch.UploadDate != nil {
		sketchView.Date = sketch.UploadDate.UTC().Format("Jan 2, 2006")
	}

	if safeDeref(sketch.Rating) != 0.0 {
		sketchView.Rating = RatingString(*sketch.Rating)
	}

	if sketch.Episode != nil && sketch.Episode.ID != nil {
		info := getShowInfo(sketch.Episode)

		sketchView.CreatorName = info.showName
		sketchView.CreatorUrl = fmt.Sprintf("/show/%d/%s", info.showId, info.showSlug)

		sketchView.CreatorImage = fmt.Sprintf("%s/show/small/%s", baseImgUrl, info.showImage)

		sketchView.CreatorInfo = fmt.Sprintf("S%dE%d · #%d", info.seNum, info.epNum, safeDeref(sketch.Number))

	} else if sketch.Creator != nil && sketch.Creator.ID != nil {
		sketchView.CreatorName = safeDeref(sketch.Creator.Name)
		sketchView.CreatorUrl = fmt.Sprintf("/creator/%d/%s", *sketch.Creator.ID, safeDeref(sketch.Creator.Slug))
		sketchView.CreatorImage = fmt.Sprintf("%s/creator/small/%s", baseImgUrl, safeDeref(sketch.Creator.ProfileImage))
	} else {
		sketchView.CreatorName = "Missing Creator"
		sketchView.CreatorImage = fmt.Sprintf("%s/missing-profile.jpg", baseImgUrl)
	}

	return sketchView, nil
}

type showInfo struct {
	epNum     int
	seNum     int
	showId    int
	showSlug  string
	showName  string
	showImage string
}

func getShowInfo(e *models.EpisodeRef) showInfo {
	info := showInfo{}
	if e == nil || e.ID == nil {
		return info
	}

	info.epNum = safeDeref(e.Number)
	if e.Season != nil {
		info.seNum = safeDeref(e.Season.Number)

		if e.Season.Show != nil {
			info.showId = safeDeref(e.Season.Show.ID)
			info.showSlug = safeDeref(e.Season.Show.Slug)
			info.showName = safeDeref(e.Season.Show.Name)
			info.showImage = safeDeref(e.Season.Show.ProfileImg)
		}
	}

	return info
}

type SketchCatalog struct {
	ResultCountLabel string
	CatalogFilter    SketchCatalogFilter
	CatalogResult    SketchCatalogResult
}

func SketchCatalogView(
	results *sketches.SketchListResult,
	htmxRequest bool,
	baseImgUrl string,
) (*SketchCatalog, error) {
	sketchCatalogResult, err := SketchCatalogResultView(
		results,
		htmxRequest,
		baseImgUrl,
	)

	fmt.Printf("%+v\n", sketchCatalogResult)
	if err != nil {
		return nil, err
	}

	sketchCatalogFilter, err := SketchCatalogFilterView(
		results,
		baseImgUrl,
	)
	if err != nil {
		return nil, err
	}

	return &SketchCatalog{
		ResultCountLabel: sketchCountLabel(results.TotalCount),
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
	results *sketches.SketchListResult,
	htmxRequest bool,
	baseImgUrl string,
) (*SketchCatalogResult, error) {
	thumbnailType := "Base"
	if len(results.Filter.PersonIDs) == 1 || len(results.Filter.CharacterIDs) == 1 {
		thumbnailType = "Cast"
	}

	sketches, err := SketchThumbnailsView(
		results.Sketches,
		baseImgUrl,
		thumbnailType,
		false,
	)
	if err != nil {
		return nil, err
	}

	currentPage := int(math.Ceil(float64(results.Filter.Offset())/float64(results.Filter.Limit()))) + 1
	totalPages := int(math.Ceil(float64(results.TotalCount) / float64(results.Filter.Limit())))
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
	if results.TotalCount != 1 {
		labelString += "es"
	}

	labelString = fmt.Sprintf(labelString, results.TotalCount)

	return &SketchCatalogResult{
		HasResults:           len(results.Sketches) != 0,
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

type SketchViewFilter struct {
	Characters []*models.Character
	Creators   []*models.Creator
	Limit      int
	Offset     int
	People     []*models.Person
	Query      string
	Shows      []*models.Show
	SortBy     string
	Tags       []*Tag
}

func SketchCatalogFilterView(result *sketches.SketchListResult, baseUrl string) (*SketchCatalogFilter, error) {
	var view SketchCatalogFilter
	sortBy := result.Filter.SortBy
	view.SortOptions = []SortOption{
		{Value: "popular", Label: "Popular", Selected: sortBy == "popular"},
		{Value: "latest", Label: "Latest", Selected: sortBy == "latest"},
		{Value: "oldest", Label: "Oldest", Selected: sortBy == "oldest"},
		{Value: "az", Label: "A-Z", Selected: sortBy == "az"},
		{Value: "za", Label: "Z-A", Selected: sortBy == "za"},
	}
	var err error
	if view.SelectedPeopleJSON, err = PeopleSelectedJSON(result.PersonRefs, baseUrl); err != nil {
		return nil, err
	}

	if view.SelectedCreatorsJSON, err = CreatorsSelectedJSON(result.CreatorRefs, baseUrl); err != nil {
		return nil, err
	}

	if view.SelectedCharactersJSON, err = CharactersSelectedJSON(result.CharacterRefs, baseUrl); err != nil {
		return nil, err
	}
	if view.SelectedShowsJSON, err = ShowsSelectedJSON(result.ShowRefs, baseUrl); err != nil {
		return nil, err
	}
	if view.SelectedTagsJSON, err = TagsSelectedJSON(result.TagRefs); err != nil {
		return nil, err
	}

	return &view, nil
}

type SelectedItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

func PeopleSelectedJSON(people []*models.PersonRef, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(people))
	for _, p := range people {

		var image string
		if p.ProfileImg != nil {
			image = fmt.Sprintf("%s/person/small/%s", baseURL, *p.ProfileImg)
		}

		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*p.ID),
			Name:  PrintPersonRefName(p),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func CreatorsSelectedJSON(creators []*models.CreatorRef, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(creators))
	for _, c := range creators {

		var image string
		if c.ProfileImage != nil {
			image = fmt.Sprintf("%s/creator/small/%s", baseURL, *c.ProfileImage)
		}
		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*c.ID),
			Name:  safeDeref(c.Name),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func ShowsSelectedJSON(shows []*models.ShowRef, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(shows))
	for _, s := range shows {

		var image string
		if s.ProfileImg != nil {
			image = fmt.Sprintf("%s/show/small/%s", baseURL, *s.ProfileImg)
		}
		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*s.ID),
			Name:  safeDeref(s.Name),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func CharactersSelectedJSON(characters []*models.CharacterRef, baseURL string) (string, error) {
	items := make([]SelectedItem, 0, len(characters))
	for _, c := range characters {

		var image string
		if c.Image != nil {
			image = fmt.Sprintf("%s/character/small/%s", baseURL, *c.Image)
		}
		items = append(items, SelectedItem{
			ID:    strconv.Itoa(*c.ID),
			Name:  safeDeref(c.Name),
			Image: image,
		})
	}

	return buildSelectedJSON(items)
}

func TagsSelectedJSON(tags []*models.TagRef) (string, error) {
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

type SketchRating struct {
	SketchID      int
	SketchTitle   string
	AverageRating string
	TotalRatings  string
	UserRating    int
}

func SketchRatingView(userSketchInfo *models.UserSketchInfo, sketch *models.Sketch) SketchRating {
	var rating int
	if userSketchInfo != nil {
		rating = safeDeref(userSketchInfo.Rating)
	}

	totalRatings := safeDeref(sketch.TotalRatings)
	var totalRatingsLabel string
	if totalRatings == 0 {
		totalRatingsLabel = "No ratings"
	} else if totalRatings == 1 {
		totalRatingsLabel += fmt.Sprintf("%d rating", totalRatings)
	} else {
		totalRatingsLabel += fmt.Sprintf("%d ratings", totalRatings)
	}
	return SketchRating{
		SketchID:      safeDeref(sketch.ID),
		SketchTitle:   safeDeref(sketch.Title),
		AverageRating: RatingString(safeDeref(sketch.Rating)),
		TotalRatings:  totalRatingsLabel,
		UserRating:    rating,
	}
}

func RatingString(rating float32) string {
	if rating == 0.0 {
		return ""
	}

	return strings.Replace(fmt.Sprintf("%.1f", rating), ".0", "", 1)
}
