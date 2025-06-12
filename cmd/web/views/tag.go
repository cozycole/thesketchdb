package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type Tag struct {
	Name string
	Url  string
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

		tagViews = append(tagViews, &Tag{Name: tagName, Url: url})
	}

	return tagViews
}
