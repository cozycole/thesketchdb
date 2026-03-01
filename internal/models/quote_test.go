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

	db := NewTestDb(t)
	sm := SketchModel{db}
	qm := QuoteModel{db}

	// insert sketch to attribute moment to
	sketchId, err := sm.Insert(&Sketch{
		Title: ptr("Test Title"),
		Slug:  ptr("test-slug"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// insert 2 quotes, one to be updated, one to be deleted on update
	q1 := &Quote{
		Text:        ptr("Quote to Update"),
		StartTimeMs: ptr(10000),
		EndTimeMs:   ptr(11000),
		Type:        ptr("quote"),
		Funny:       ptr("funny"),
	}

	q2 := &Quote{
		Text:        ptr("Quote to Delete"),
		StartTimeMs: ptr(12000),
		EndTimeMs:   ptr(13000),
		Type:        ptr("action"),
		Funny:       ptr("support"),
	}

	ctx := context.Background()
	tx, err := qm.DB.Begin(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	// insert setup goes here
	q1Id, err := insertQuote(ctx, tx, q1, sketchId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = insertQuote(ctx, tx, q2, sketchId)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	// END OF SETUP

	quotesForUpdate := []*Quote{
		{
			Text:        ptr("Newly inserted quote"),
			StartTimeMs: ptr(15000),
			EndTimeMs:   ptr(20000),
		},
		{
			ID:          ptr(q1Id),
			Text:        ptr("Quote has been updated"),
			StartTimeMs: ptr(20000),
			EndTimeMs:   ptr(25000),
			Funny:       ptr("support"),
		},
	}

	err = qm.BatchUpdateQuotes(sketchId, quotesForUpdate)
	if err != nil {
		t.Fatal(err)
	}

	updatedQuotes, err := qm.GetBySketch(sketchId)
	if err != nil {
		t.Fatal(err)
	}

	if len(updatedQuotes) == 0 {
		t.Fatal("no quotes found to compare")
	}

	if len(updatedQuotes) != 2 {
		assert.Equal(t, len(updatedQuotes), 2)
		return
	}

	// assert first quote is new
	assert.Equal(t, *updatedQuotes[0].ID, 3)
	assert.Equal(t, *updatedQuotes[0].Text, "Newly inserted quote")
	assert.Equal(t, *updatedQuotes[0].StartTimeMs, 15000)
	assert.Equal(t, *updatedQuotes[0].EndTimeMs, 20000)

	// assert second quote is updated correctly
	assert.Equal(t, *updatedQuotes[1].ID, 1)
	assert.Equal(t, *updatedQuotes[1].Text, "Quote has been updated")
	assert.Equal(t, *updatedQuotes[1].StartTimeMs, 20000)
	assert.Equal(t, *updatedQuotes[1].EndTimeMs, 25000)
	assert.Equal(t, *updatedQuotes[1].Funny, "support")
	assert.Equal(t, updatedQuotes[1].Type, nil)
}

func TestQuoteTagUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	db := NewTestDb(t)
	sm := SketchModel{db}
	qm := QuoteModel{db}
	tm := TagModel{db}

	// insert sketch to attribute moment to
	sketchId, err := sm.Insert(&Sketch{
		Title: ptr("Test Title"),
		Slug:  ptr("test-slug"),
	})
	if err != nil {
		t.Fatal(err)
	}

	q1 := &Quote{
		Text:        ptr("Quote to Update"),
		StartTimeMs: ptr(10000),
		EndTimeMs:   ptr(11000),
		Type:        ptr("quote"),
		Funny:       ptr("funny"),
	}

	err = qm.BatchUpdateQuotes(sketchId, []*Quote{q1})
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
	tx, err := qm.DB.Begin(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	quoteId := 1

	err = insertQuoteTagAssociations(ctx, tx, quoteId, []int{*t1.ID, *t2.ID})

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	// tag id 1 and 2 are associated with quote one
	// we want id 2 and 3 to remain after update
	qm.BatchUpdateQuoteTags(1, []int{*t2.ID, *t3.ID})

	quotes, err := qm.GetBySketch(sketchId)
	if err != nil {
		t.Fatal(err)
	}

	if len(quotes) != 1 {
		t.Fatal(fmt.Errorf("a single quote is expected"))
	}

	desiredIds := []int{2, 3}
	for _, tag := range quotes[0].Tags {
		if !slices.Contains(desiredIds, *tag.ID) {
			t.Errorf("tag id %d not a desired tag", *tag.ID)
		}
	}
}

func TestQuoteCastUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	db := NewTestDb(t)
	sm := SketchModel{db}
	qm := QuoteModel{db}
	cm := CastModel{db}

	// insert sketch to attribute moment to
	sketchId, err := sm.Insert(&Sketch{
		Title: ptr("Test Title"),
		Slug:  ptr("test-slug"),
	})
	if err != nil {
		t.Fatal(err)
	}

	q1 := &Quote{
		Text:        ptr("Quote to Update"),
		StartTimeMs: ptr(10000),
		EndTimeMs:   ptr(11000),
		Type:        ptr("quote"),
		Funny:       ptr("funny"),
	}

	err = qm.BatchUpdateQuotes(sketchId, []*Quote{q1})
	if err != nil {
		t.Fatal(err)
	}

	cm1 := &CastMember{
		CharacterName: ptr("Character Name #1"),
	}

	cm2 := &CastMember{
		CharacterName: ptr("Character Name #2"),
	}

	cm3 := &CastMember{
		CharacterName: ptr("Character Name #3"),
	}

	cm.Insert(sketchId, cm1)
	cm.Insert(sketchId, cm2)
	cm.Insert(sketchId, cm3)

	ctx := context.Background()
	tx, err := qm.DB.Begin(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	quoteId := 1

	err = insertQuoteCastAssociations(ctx, tx, quoteId, []int{*cm1.ID, *cm2.ID})

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	// cast id 1 and 2 are associated with quote one
	// we want id 2 and 3 to remain after update
	qm.BatchUpdateQuoteCastMembers(1, []int{*cm2.ID, *cm3.ID})

	quotes, err := qm.GetBySketch(sketchId)
	if err != nil {
		t.Fatal(err)
	}

	if len(quotes) != 1 {
		t.Fatal(fmt.Errorf("a single quote is expected"))
	}

	desiredIds := []int{*cm2.ID, *cm3.ID}
	for _, cast := range quotes[0].CastMembers {
		if !slices.Contains(desiredIds, *cast.ID) {
			t.Errorf("tag id %d not a desired tag", *cast.ID)
		}
	}
}
