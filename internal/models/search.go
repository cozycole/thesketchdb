package models

import ()

// NOTE: Query is defined on the Filter, SearchResult and templateData structs
// Given the search term: kenan snl
// - Filter.Query -> "kenan | snl"
// - SearchResult.Query -> "kenan+snl"
// - templateData.Query -> "kenan snl" (i.e. user facing)

type SearchResult struct {
	Type                string
	Query               string
	SketchResults       []*Sketch
	TotalSketchCount    int
	PersonResults       []*Person
	TotalPersonCount    int
	CreatorResults      []*Creator
	TotalCreatorCount   int
	CharacterResults    []*Character
	TotalCharacterCount int
	ShowResults         []*Show
	TotalShowCount      int
	Filter              *Filter
	CurrentPage         int
	Pages               []int
}
