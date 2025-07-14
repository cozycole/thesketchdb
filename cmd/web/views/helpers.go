package views

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
	"unicode"

	"sketchdb.cozycole.net/internal/models"
)

func printCast(cast []*models.CastMember) string {
	castList := ""
	var personIds []int
	for i, cm := range cast {
		if cm.Actor == nil ||
			cm.Actor.ID == nil ||
			intSliceContains(personIds, *cm.Actor.ID) {
			continue
		}

		name := PrintPersonName(cm.Actor)
		if name != "" {
			if i != 0 {
				name = ", " + name
			}
			castList += name
		}

		personIds = append(personIds, *cm.Actor.ID)
	}

	return castList
}

func PrintPersonName(a *models.Person) string {
	if a == nil {
		return ""
	}
	var name string
	if a.First != nil {
		name = *a.First
	}

	if a.Last == nil {
		return name
	}

	return name + " " + *a.Last
}

func PrintEpisodeName(e *models.Episode) string {
	if e == nil {
		return ""
	}

	out := safeDeref(e.ShowName)
	if e.SeasonNumber != nil {
		out += fmt.Sprintf(" S%d", safeDeref(e.SeasonNumber))
	}
	if e.Number != nil {
		out += fmt.Sprintf("E%d", safeDeref(e.Number))
	}

	return out
}

func uppercaseFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func humanDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format("Jan 2, 2006")
}

func createEpisodeTitle(episode *models.Episode) string {
	title := safeDeref(episode.Title)
	if title == "" {
		if episode.Number == nil {
			title = "Episode ?"
		} else {
			title = fmt.Sprintf("Episode %d", *episode.Number)
		}
	}
	return title
}

func determineEpisodeWatchURL(episode *models.Episode) (string, string) {
	if safeDeref(episode.YoutubeID) != "" {
		return fmt.Sprintf("https://www.youtube.com/watch?v=%s", *episode.YoutubeID),
			"/static/img/youtube-logo.jpg"
	}
	return "", ""
}

func seasonEpisodeInfo(episode *models.Episode) string {
	var info string

	if episode.SeasonNumber != nil &&
		episode.Number != nil {

		info = fmt.Sprintf(
			"S%d · E%d · %s",
			*episode.SeasonNumber,
			*episode.Number,
			sketchCountLabel(len(episode.Sketches)),
		)
	}

	return info
}

func seasonEpisodeUrl(episode *models.Episode) string {
	var url string
	if episode.ShowID != nil &&
		episode.ShowSlug != nil &&
		episode.SeasonNumber != nil &&
		episode.Number != nil {

		url = fmt.Sprintf(
			"/show/%d/%s/season/%d/episode/%d",
			*episode.ShowID,
			*episode.ShowSlug,
			*episode.SeasonNumber,
			*episode.Number,
		)
	}

	return url
}

func getAge(birthDate *time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()

	if today.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

func BuildURL(baseURL string, currentPage int, filter *models.Filter) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	params := filter.Params()
	params.Add("page", strconv.Itoa(currentPage))
	if filter.Query != "" {
		params.Set("query", filter.Query)
	}

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func safeDeref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

func sketchCountLabel(count int) string {
	labelString := "%d Sketch"
	if count != 1 {
		labelString += "es"
	}

	return fmt.Sprintf(labelString, count)
}

func intSliceContains(list []int, value int) bool {
	for _, n := range list {
		if n == value {
			return true
		}
	}
	return false
}

func getSeasonSketchCount(season *models.Season) int {
	var count int
	for _, ep := range season.Episodes {
		count += len(ep.Sketches)
	}

	return count
}

func getShowEpisodeCount(show *models.Show) int {
	var count int
	for _, season := range show.Seasons {
		count += len(season.Episodes)
	}

	return count
}

func getShowSketchCount(show *models.Show) int {
	var count int
	for _, season := range show.Seasons {
		for _, ep := range season.Episodes {
			count += len(ep.Sketches)
		}
	}

	return count
}
