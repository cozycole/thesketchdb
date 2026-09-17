package models

import (
	"math"
	"net/url"
	"strconv"
)

type Filter struct {
	Page         int
	PageSize     int // Page is 1 indexed
	Query        string
	Type         string
	CharacterIDs []int
	CreatorIDs   []int
	PersonIDs    []int
	SketchIDs    []int
	ShowIDs      []int
	TagIDs       []int
	SortBy       string
}

func (f Filter) Limit() int {
	if f.PageSize < 1 {
		return 24
	}
	return f.PageSize
}

func (f Filter) Offset() int {
	if f.Page < 1 {
		return 0
	}

	return (f.Page - 1) * f.PageSize
}

var sortMap = map[string]string{
	"popular": "popularity DESC, upload_date DESC",
	"recent":  "sketch_id DESC",
	"newest":  "upload_date DESC, sketch_title ASC",
	"oldest":  "upload_date ASC, sketch_title ASC",
	"az":      "sketch_title ASC",
	"za":      "sketch_title DESC",
}

func (f *Filter) Params() url.Values {
	params := url.Values{}

	if f.SortBy != "" {
		params.Add("sort", f.SortBy)
	}

	if f.Query != "" {
		params.Add("query", url.QueryEscape(f.Query))
	}

	for _, id := range f.PersonIDs {
		params.Add("person", strconv.Itoa(id))
	}

	for _, id := range f.CreatorIDs {
		params.Add("creator", strconv.Itoa(id))
	}

	for _, id := range f.ShowIDs {
		params.Add("show", strconv.Itoa(id))
	}

	for _, id := range f.CharacterIDs {
		params.Add("character", strconv.Itoa(id))
	}

	for _, id := range f.TagIDs {
		params.Add("tag", strconv.Itoa(id))
	}

	return params
}

func (f *Filter) ParamsString() string {
	return f.Params().Encode()
}

type Metadata struct {
	CurrentPage  int `json:"page"`
	PageSize     int `json:"pageSize"`
	TotalPages   int `json:"totalPages"`
	TotalRecords int `json:"total"`
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
