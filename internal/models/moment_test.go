package models

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

func TestBatchQuoteUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	db := newTestDB(t)
	sm := SketchModel{db}
	cm := CastModel{db}
	mm := MomentModel{db}

	// insert sketch to attribute moment to
	sketchId, err := sm.Insert(&Sketch{
		Title: ptr("Test Title"),
		Slug:  ptr("test-slug"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// add one cast member to attribute quote to
	// and another to test update
	castId, err := cm.Insert(sketchId, &CastMember{CharacterName: ptr("Test Character")})
	if err != nil {
		t.Fatal(err)
	}

	cast2Id, err := cm.Insert(sketchId, &CastMember{CharacterName: ptr("Test Character 2")})
	if err != nil {
		t.Fatal(err)
	}

	// insert moment to update
	momentId, err := mm.Insert(sketchId, &Moment{Timestamp: ptr(1)})
	if err != nil {
		t.Fatal(err)
	}

	// insert 2 quotes, one to be updated, one to be deleted on update
	q1 := &Quote{
		CastMember: &CastMember{ID: ptr(castId)},
		Text:       ptr("Quote to Update"),
		Position:   ptr(1),
	}

	q2 := &Quote{
		CastMember: &CastMember{ID: ptr(castId)},
		Text:       ptr("Quote to Delete"),
		Position:   ptr(2),
	}
	ctx := context.Background()
	tx, err := mm.DB.Begin(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	q1Id, err := insertQuote(ctx, tx, q1, momentId)
	_, err2 := insertQuote(ctx, tx, q2, momentId)
	if err != nil || err2 != nil {
		t.Fatal(fmt.Errorf("error inserting quotes on setup: %w", err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	// END OF SETUP

	quotes := []*Quote{
		{
			CastMember: q1.CastMember,
			Text:       ptr("Newly inserted quote"),
			Position:   ptr(1),
		},
		{
			ID:         ptr(q1Id),
			CastMember: &CastMember{ID: ptr(cast2Id)},
			Text:       ptr("Quote has been updated"),
			Position:   ptr(2),
		},
	}

	err = mm.BatchUpdateQuotes(momentId, quotes)
	if err != nil {
		t.Fatal(err)
	}

	moments, err := mm.GetBySketch(sketchId)
	if err != nil {
		t.Fatal(err)
	}

	if len(moments) == 0 {
		t.Fatal("no moments found to compare")
	}

	updatedQuotes := moments[0].Quotes

	if len(updatedQuotes) != 2 {
		assert.Equal(t, len(updatedQuotes), 2)
		return
	}

	// assert first quote is new
	assert.Equal(t, *updatedQuotes[0].Text, "Newly inserted quote")
	assert.Equal(t, *updatedQuotes[0].Position, 1)

	// assert second quote is updated correctly
	assert.Equal(t, *updatedQuotes[1].ID, 1)
	assert.Equal(t, *updatedQuotes[1].Text, "Quote has been updated")
	assert.Equal(t, *updatedQuotes[1].Position, 2)
	assert.Equal(t, *updatedQuotes[1].CastMember.ID, cast2Id)
}

func TestQuoteTagUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	db := newTestDB(t)
	sm := SketchModel{db}
	cm := CastModel{db}
	mm := MomentModel{db}
	tm := TagModel{db}

	// insert sketch to attribute moment to
	sketchId, err := sm.Insert(&Sketch{
		Title: ptr("Test Title"),
		Slug:  ptr("test-slug"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// add one cast member to attribute quote to
	// and another to test update
	castId, err := cm.Insert(sketchId, &CastMember{CharacterName: ptr("Test Character")})
	if err != nil {
		t.Fatal(err)
	}

	// insert moment to update
	momentId, err := mm.Insert(sketchId, &Moment{Timestamp: ptr(1)})
	if err != nil {
		t.Fatal(err)
	}

	// insert 2 quotes, one to be updated, one to be deleted on update
	q1 := &Quote{
		CastMember: &CastMember{ID: ptr(castId)},
		Text:       ptr("Quote to Update"),
		Position:   ptr(1),
	}

	err = mm.BatchUpdateQuotes(momentId, []*Quote{q1})
	if err != nil {
		t.Fatal(err)
	}

	t1 := &Tag{
		Name: ptr("Tag 1"),
		Slug: ptr("tag-1"),
		Type: ptr("quote"),
	}

	t2 := &Tag{
		Name: ptr("Tag 2"),
		Slug: ptr("tag-2"),
		Type: ptr("quote"),
	}

	t3 := &Tag{
		Name: ptr("Tag 3"),
		Slug: ptr("tag-3"),
		Type: ptr("quote"),
	}

	tm.Insert(t1)
	tm.Insert(t2)
	tm.Insert(t3)

	ctx := context.Background()
	tx, err := mm.DB.Begin(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	quoteId := 1

	err = insertQuoteTagAssociations(ctx, tx, quoteId, []int{1, 2})

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	// tag id 1 and 2 are associated with quote one
	// we want id 2 and 3 to remain after update
	mm.BatchUpdateQuoteTags(1, []*Tag{t2, t3})

	quote, err := mm.GetQuote(quoteId)
	if err != nil {
		t.Fatal(err)
	}

	desiredIds := []int{2, 3}
	for _, tag := range quote.Tags {
		if !slices.Contains(desiredIds, *tag.ID) {
			t.Errorf("tag id %d not a desired tag", *tag.ID)
		}
	}
}
