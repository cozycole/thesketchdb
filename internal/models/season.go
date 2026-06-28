package models

import (
	"strconv"
	"time"
)

type Season struct {
	ID       *int
	Slug     *string
	Number   *int
	Show     *ShowRef
	Episodes []*EpisodeRef
}

type SeasonRef struct {
	ID     *int     `json:"id"`
	Slug   *string  `json:"slug"`
	Number *int     `json:"number"`
	Show   *ShowRef `json:"show"`
}

func (s *Season) AirYear() string {
	var airDates []time.Time
	for _, e := range s.Episodes {
		if e.AirDate != nil {
			airDates = append(airDates, *e.AirDate)
		}
	}
	if len(airDates) == 0 {
		return ""
	}

	min := airDates[0]
	for _, t := range airDates[1:] {
		if t.Before(min) {
			min = t
		}
	}
	return strconv.Itoa(min.Year())
}
