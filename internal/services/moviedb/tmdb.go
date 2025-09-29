package moviedb

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseIMDbID extracts the IMDb person ID from a full URL
func ParseIMDbID(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "name" {
		return "", fmt.Errorf("invalid IMDb person URL: %s", rawURL)
	}
	return parts[1], nil
}

// ParseTMDBID extracts the TMDB person ID from a full URL
func ParseTMDbID(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "person" {
		return "", fmt.Errorf("invalid TMDB person URL: %s", rawURL)
	}
	return parts[1], nil
}

func BuildIMDbURL(id string) string {
	return "https://www.imdb.com/name/" + id
}

func BuildTMDbURL(id string) string {
	return "https://www.themoviedb.org/person/" + id
}
