package views

import (
	"fmt"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

type Quote struct {
	ID                    int
	StartTimestamp        string
	StartTimestampSeconds int
	Text                  string
	CastLabel             string
	CastImgUrls           []string
	ExtraCast             int
	IsLiked               bool
	LikeCount             int
	InsertDivider         bool
}

const MAX_DISPLAY_IMAGES = 4
const QUOTE_TIMESTAMP_MS_DIVIDER = 10000

func SketchQuoteSection(quotes []*models.Quote, baseImgUrl string) []Quote {
	viewQuotes := []Quote{}
	previousTimestamp := 0
	for i, q := range quotes {
		viewQuote := Quote{}

		timestamp := safeDeref(q.StartTimeMs)
		viewQuote.ID = safeDeref(q.ID)
		viewQuote.StartTimestamp = models.MillisecondsToMMSS(timestamp)
		viewQuote.StartTimestampSeconds = timestamp
		viewQuote.Text = safeDeref(q.Text)
		viewQuote.CastLabel = QuoteHeader(q.CastMembers)
		viewQuote.IsLiked = safeDeref(q.UserLiked)
		viewQuote.LikeCount = safeDeref(q.LikeCount)

		// logic to insert dividers between long pauses between quotes
		if i != 0 && safeDeref(q.StartTimeMs)-previousTimestamp > QUOTE_TIMESTAMP_MS_DIVIDER {
			viewQuote.InsertDivider = true
		}
		previousTimestamp = safeDeref(q.StartTimeMs)
		if q.EndTimeMs != nil {
			previousTimestamp = *q.EndTimeMs
		}

		for _, cm := range q.CastMembers {
			viewQuote.CastImgUrls = append(
				viewQuote.CastImgUrls,
				DetermineCastImageUrl(cm, "small", baseImgUrl),
			)
		}
		viewQuote.ExtraCast = max(len(viewQuote.CastImgUrls)-MAX_DISPLAY_IMAGES, 0)
		viewQuotes = append(viewQuotes, viewQuote)
	}

	return viewQuotes
}

func QuoteHeader(members []*models.CastMember) string {
	if len(members) == 0 {
		return ""
	}

	if len(members) == 1 {
		cm := members[0]
		charName := safeDeref(cm.CharacterName)

		actorName := PrintPersonRefName(cm.Actor)
		if actorName != "" {
			if charName != "" {
				return charName + fmt.Sprintf(" (%s) ", PrintPersonRefName(cm.Actor))
			}
			return actorName
		}
	}

	// for just two members we keep the actor's names in parenthesis
	if len(members) == 2 {
		names := []string{}
		for _, cm := range members {
			charName := safeDeref(cm.CharacterName)

			actorName := PrintPersonRefName(cm.Actor)
			if actorName != "" {
				charName = charName + fmt.Sprintf(" (%s)", PrintPersonRefName(cm.Actor))
			}
			names = append(names, charName)
		}

		return strings.Join(names, ", ")
	}

	charNames := []string{}
	for _, cm := range members {
		charName := safeDeref(cm.CharacterName)
		if charName != "" {
			charNames = append(charNames, charName)
		}
	}

	return strings.Join(charNames, ", ")
}

type QuoteListItem struct {
	ID                   int
	Text                 string
	CastMembers          []QuoteCastMember
	ExtraCastImagesCount int
	ExtraCastNamesCount  int
	IsLiked              bool
	LikeCount            int

	SketchURL   string
	SketchTitle string

	Show QuoteLinkRef

	Creators []QuoteLinkRef
}

type QuoteLinkRef struct {
	URL   string
	Title string
}

type QuoteCastMember struct {
	PersonName    string
	PersonURL     string
	CharacterName string
	CharacterURL  string
	ImageURL      string
}

const MAX_DISPLAY_NAMES = 2

func QuoteListView(quotes []*models.Quote, baseImgUrl string) []QuoteListItem {
	quoteItems := []QuoteListItem{}
	for _, q := range quotes {
		qi := QuoteListItem{}
		qi.ID = safeDeref(q.ID)

		qi.Text = safeDeref(q.Text)
		for _, cm := range q.CastMembers {
			qcm := QuoteCastMember{
				ImageURL: DetermineCastImageUrl(cm, "small", baseImgUrl),
			}
			if cm.Actor != nil {
				qcm.PersonName = PrintPersonRefName(cm.Actor)
				qcm.PersonURL = fmt.Sprintf(
					"/person/%d/%s",
					safeDeref(cm.Actor.ID),
					safeDeref(cm.Actor.Slug),
				)
			}

			qcm.CharacterName = safeDeref(cm.CharacterName)
			if cm.Character != nil {
				if qcm.CharacterName == "" {
					qcm.CharacterName = safeDeref(cm.Character.Name)
				}
				qcm.PersonName = PrintPersonRefName(cm.Actor)
				qcm.CharacterURL = fmt.Sprintf(
					"/character/%d/%s",
					safeDeref(cm.Character.ID),
					safeDeref(cm.Character.Slug),
				)
			}
			qi.CastMembers = append(qi.CastMembers, qcm)
		}
		// we want to show up to 4 cast images and up to 2 cast names
		qi.ExtraCastImagesCount = max(len(qi.CastMembers)-MAX_DISPLAY_IMAGES, 0)
		qi.ExtraCastNamesCount = max(len(qi.CastMembers)-MAX_DISPLAY_NAMES, 0)
		qi.LikeCount = safeDeref(q.LikeCount)

		qi.SketchURL = fmt.Sprintf(
			"/sketch/%d/%s",
			safeDeref(q.SketchID),
			safeDeref(q.SketchSlug),
		)

		qi.IsLiked = safeDeref(q.UserLiked)
		qi.SketchTitle = safeDeref(q.SketchTitle)
		if q.Show != nil {
			qi.Show.URL = fmt.Sprintf(
				"/show/%d/%s",
				safeDeref(q.Show.ID),
				safeDeref(q.Show.Slug),
			)
			qi.Show.Title = safeDeref(q.Show.Name)
		}

		for _, c := range q.Creators {
			qi.Creators = append(qi.Creators, QuoteLinkRef{
				URL: fmt.Sprintf(
					"/creator/%d/%s",
					safeDeref(c.ID),
					safeDeref(c.Slug),
				),
				Title: safeDeref(c.Name),
			})
		}

		quoteItems = append(quoteItems, qi)
	}

	return quoteItems
}
