package sketches

import (
	"fmt"
	"net/url"

	"sketchdb.cozycole.net/internal/models"
)

var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func createSketchSlug(sketch *models.Sketch) string {
	var slugInput string

	if sketch.Episode != nil {
		ep := sketch.Episode
		show := ep.GetShow()
		if show != nil {
			seasonNumber := safeDeref(ep.Season.Number)
			showString := safeDeref(show.Name)
			episodeNumber := safeDeref(ep.Number)
			slugInput += fmt.Sprintf("%s s%d e%d", showString, seasonNumber, episodeNumber)
		}
	}

	if sketch.Creator != nil {
		slugInput += safeDeref(sketch.Creator.Name)
	}

	if slugInput == "" {
		return safeDeref(sketch.Title)
	}

	return models.CreateSlugName(slugInput + " " + safeDeref(sketch.Title))
}

func extractYouTubeVideoID(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	if videoID == "" {
		return "", fmt.Errorf("video ID not found in URL")
	}

	return videoID, nil
}

func safeDeref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}
