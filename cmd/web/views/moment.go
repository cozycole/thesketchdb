package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type Moment struct {
	MomentID       int
	QuoteTableForm QuoteTableForm
}

func UpdateMomentView(moment models.Moment, baseImgUrl string) Moment {
	return Moment{
		MomentID:       safeDeref(moment.ID),
		QuoteTableForm: QuoteTableFormView(moment, baseImgUrl),
	}
}

type QuoteTableForm struct {
	MomentID  int
	QuoteRows []QuoteRow
}

type QuoteRow struct {
	ID        int
	ImageUrl  string
	ActorName string
	LineText  string
	Type      string
	Funny     string
	AddTagUrl string
	TagCount  int
}

func QuoteTableFormView(moment models.Moment, baseImgUrl string) QuoteTableForm {
	quoteTable := QuoteTableForm{}
	quoteTable.MomentID = safeDeref(moment.ID)
	for _, q := range moment.Quotes {
		row := QuoteRow{}
		row.ID = safeDeref(q.ID)
		cm := safeDeref(q.CastMember)
		if cm.ProfileImg != nil {
			row.ImageUrl = fmt.Sprintf("%s/cast/profile/small/%s", baseImgUrl, safeDeref(cm.ProfileImg))
		} else if cm.Actor != nil {
			row.ImageUrl = fmt.Sprintf("%s/person/small/%s", baseImgUrl, safeDeref(cm.Actor.ProfileImg))
		} else {
			row.ImageUrl = "/static/img/missing-profile.jpg"
		}

		if cm.Actor != nil {
			row.ActorName = PrintPersonName(cm.Actor)
			row.ActorName = fmt.Sprintf(
				"/person/%d/%s",
				safeDeref(cm.Actor.ID),
				safeDeref(cm.Actor.Slug),
			)
		}

		row.Type = UppercaseFirst(safeDeref(q.Type))
		row.Funny = UppercaseFirst(safeDeref(q.Funny))

		quoteTable.QuoteRows = append(quoteTable.QuoteRows, row)
	}

	return quoteTable
}
