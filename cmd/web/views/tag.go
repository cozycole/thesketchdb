package views

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"sketchdb.cozycole.net/internal/models"
)

type Tag struct {
	ID       int
	Name     string
	Url      string
	Selected bool
}

func TagsView(tags []*models.Tag) []*Tag {
	var tagViews []*Tag
	for _, tag := range tags {
		var tagName, url string

		if tag.ID != nil && tag.Name != nil {
			tagName = *tag.Name
		}

		if tag.Category.ID != nil && tag.Category.Name != nil {
			tagName = fmt.Sprintf("%s / %s", *tag.Category.Name, tagName)
		}

		if tag.ID != nil {
			url = fmt.Sprintf("/catalog/sketches?tag=%d", *tag.ID)
		}

		tagName = TitleCaseEnglish(tagName)
		tagViews = append(tagViews, &Tag{ID: safeDeref(tag.ID), Name: tagName, Url: url})
	}

	return tagViews
}

type TagTable struct {
	SketchID int
	Error    string
	TagRows  []TagRow
}

type TagRow struct {
	ID   int
	Name string
}

func TagTableView(tags []*models.Tag, sketchID int) TagTable {
	rows := make([]TagRow, 0, len(tags))
	for _, tag := range tags {
		row := TagRow{}
		row.ID = safeDeref(tag.ID)

		var tagName string
		if tag.Category != nil && safeDeref(tag.Category.Name) != "" {
			tagName = safeDeref(tag.Category.Name) + " / "
		}

		tagName += safeDeref(tag.Name)

		c := cases.Title(language.English)
		row.Name = c.String(tagName)
		rows = append(rows, row)
	}
	return TagTable{
		SketchID: sketchID,
		TagRows:  rows,
	}
}
