package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type UserPage struct {
	Username          string
	DateJoined        string
	Image             string
	FavoritedSketches *SketchGallery
}

func UserPageView(user *models.User, favorited []*models.Sketch, baseImgUrl string) (*UserPage, error) {
	if user.ID == nil {
		return nil, fmt.Errorf("User ID not defined")
	}

	page := UserPage{}

	page.Username = safeDerefString(user.Username)
	page.DateJoined = humanDate(user.CreatedAt)
	page.Image = "/static/img/missing-profile.jpg"
	if user.ProfileImage != nil {
		page.Image = fmt.Sprintf(
			"%s/user/%s",
			baseImgUrl,
			*user.ProfileImage,
		)
	}

	var err error
	page.FavoritedSketches, err = SketchGalleryView(
		favorited,
		baseImgUrl,
		"base",
		"sub",
		12,
	)
	if err != nil {
		return nil, err
	}

	return &page, nil
}
