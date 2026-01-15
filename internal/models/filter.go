package models

import (
	"net/url"
	"strconv"
)

type Filter struct {
	Limit        int
	Offset       int
	Query        string
	CharacterIDs []int
	CreatorIDs   []int
	PersonIDs    []int
	ShowIDs      []int
	TagIDs       []int
	SortBy       string
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
