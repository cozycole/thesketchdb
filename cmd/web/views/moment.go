package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type Moment struct {
	ID        int
	Timestamp string
	Quotes    []Quote
}

type Quote struct {
	CastName   string
	CastImgUrl string
	Text       string
}

func SketchQuoteSection(moments []*models.Moment, baseImgUrl string) []Moment {
	viewMoments := []Moment{}
	for _, m := range moments {
		vm := Moment{}
		vm.ID = safeDeref(m.ID)
		vm.Timestamp = models.SecondsToMMSS(safeDeref(m.Timestamp))

		for _, q := range m.Quotes {
			vq := Quote{}
			vq.CastName = QuoteHeader(q.CastMember)
			vq.CastImgUrl = DetermineCastImageUrl(q.CastMember, "small", baseImgUrl)
			vq.Text = safeDeref(q.Text)

			vm.Quotes = append(vm.Quotes, vq)
		}

		viewMoments = append(viewMoments, vm)
	}

	return viewMoments
}

func QuoteHeader(cm *models.CastMember) string {
	if cm == nil {
		return ""
	}

	charName := safeDeref(cm.CharacterName)

	actorName := PrintPersonName(cm.Actor)
	if actorName != "" {
		if charName != "" {
			return charName + fmt.Sprintf(" (%s)", PrintPersonName(cm.Actor))
		}
		return actorName
	}

	return charName
}
