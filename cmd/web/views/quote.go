package views

import (
	"fmt"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

type Quote struct {
	ID             int
	StartTimestamp string
	Text           string
	CastLabel      string
	CastImgUrls    []string
	ExtraCast      int
	InsertDivider  bool
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
		viewQuote.Text = safeDeref(q.Text)
		viewQuote.CastLabel = QuoteHeader(q.CastMembers)

		// logic to insert dividers between long pauses between quotes
		fmt.Printf("Previous: %d Current: %d\n", previousTimestamp, timestamp)
		if i != 0 && safeDeref(q.StartTimeMs)-previousTimestamp > QUOTE_TIMESTAMP_MS_DIVIDER {
			viewQuote.InsertDivider = true
		}
		previousTimestamp = safeDeref(q.StartTimeMs)

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
