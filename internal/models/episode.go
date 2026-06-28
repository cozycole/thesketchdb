package models

import (
	"time"
)

type Episode struct {
	ID        *int         `json:"id"`
	Slug      *string      `json:"slug"`
	Title     *string      `json:"title"`
	Number    *int         `json:"number"`
	AirDate   *time.Time   `json:"airDate"`
	Thumbnail *string      `json:"thumbnail"`
	Season    *SeasonRef   `json:"season"`
	URL       *string      `json:"url"`
	Sketches  []*SketchRef `json:"sketches"`
	YoutubeID *string      `json:"youtubeId"`
}

func (e *Episode) GetTitle() *string     { return e.Title }
func (e *Episode) GetNumber() *int       { return e.Number }
func (e *Episode) GetSeason() *SeasonRef { return e.Season }
func (e *Episode) GetSketchCount() int   { return len(e.Sketches) }

func (e *Episode) GetShow() *ShowRef {
	if e.Season == nil ||
		e.Season.ID == nil ||
		e.Season.Show == nil ||
		e.Season.Show.ID == nil {
		return nil
	}
	return e.Season.Show
}

func (e *Episode) ToRef() *EpisodeRef {
	if e.ID == nil {
		return nil
	}

	sketchCount := len(e.Sketches)
	return &EpisodeRef{
		ID:          e.ID,
		Slug:        e.Slug,
		Title:       e.Title,
		Number:      e.Number,
		AirDate:     e.AirDate,
		Thumbnail:   e.Thumbnail,
		Season:      e.Season,
		SketchCount: &sketchCount,
	}
}

type EpisodeRef struct {
	ID          *int       `json:"id"`
	Slug        *string    `json:"slug"`
	Title       *string    `json:"title"`
	Number      *int       `json:"number"`
	AirDate     *time.Time `json:"airDate"`
	Thumbnail   *string    `json:"thumbnail"`
	Season      *SeasonRef `json:"season"`
	SketchCount *int       `json:"-"`
}

func (e *EpisodeRef) GetTitle() *string     { return e.Title }
func (e *EpisodeRef) GetNumber() *int       { return e.Number }
func (e *EpisodeRef) GetSeason() *SeasonRef { return e.Season }
func (e *EpisodeRef) GetSketchCount() int   { return safeDeref(e.SketchCount) }

func (e *EpisodeRef) GetShow() *ShowRef {
	if e.Season == nil ||
		e.Season.ID == nil ||
		e.Season.Show.ID == nil {
		return nil
	}
	return e.Season.Show
}
