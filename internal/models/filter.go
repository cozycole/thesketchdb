package models

import (
	"math"
	"net/url"
	"strconv"
)

type Filter struct {
	Page         int
	PageSize     int
	Query        string
	CharacterIDs []int
	CreatorIDs   []int
	PersonIDs    []int
	ShowIDs      []int
	TagIDs       []int
	SortBy       string
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

var sortMap = map[string]string{
	"popular": "popularity DESC, upload_date DESC",
	"latest":  "upload_date DESC, sketch_title ASC",
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
// values given the total number of records, current page, and page size values. Note
// that when the last page value is calculated we are dividing two int values, and
// when dividing integer types in Go the result will also be an integer type, with
// the modulus (or remainder) dropped. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be (12+5-1)/5 = 3.2, which is then
// truncated to 3 by Go.
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
