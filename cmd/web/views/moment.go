package views

import (
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
			vq.CastName = safeDeref(q.CastMember.CharacterName)
			vq.CastImgUrl = determineCastImageUrl(q.CastMember, "small", baseImgUrl)
			vq.Text = safeDeref(q.Text)

			vm.Quotes = append(vm.Quotes, vq)
		}

		viewMoments = append(viewMoments, vm)
	}

	return viewMoments
}
