package casts

import (
	"fmt"
	"sketchdb.cozycole.net/internal/models"
	"slices"
)

func validateCastIds(castIds []int, castMembers []*models.CastMember) error {
	allowed := make(map[int]struct{}, len(castMembers))
	for _, cm := range castMembers {
		if cm != nil && cm.ID != nil {
			allowed[*cm.ID] = struct{}{}
		}
	}
	for _, id := range castIds {
		if _, ok := allowed[id]; !ok {
			return fmt.Errorf("invalid cast id: %d", id)
		}
	}
	return nil
}

func getNextPosition(castMembers []*models.CastMember) int {
	positions := []int{}
	for _, cm := range castMembers {
		positions = append(positions, *cm.Position)
	}

	if len(positions) == 0 {
		return 1
	}

	return slices.Max(positions) + 1
}
