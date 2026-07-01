package main

import (
	"fmt"
	"slices"
)

type PostAction interface {
	ApplyFilter(rolls []Roll) []Roll
}

type KeepHighest struct {
	keepCount int
}

func (postAction KeepHighest) String() string {
	return fmt.Sprintf("kh%d", postAction.keepCount)
}

func (khAction KeepHighest) ApplyFilter(rolls []Roll) []Roll {
	indices := make([]int, len(rolls))
	for i := range indices {
		indices[i] = i
	}

	// Sort a slice of indices by their corresponding roll values (descending).
	// This lets us rank rolls by value without reordering the original slice,
	// so we can mark which positions are dropped while preserving insertion order.
	slices.SortFunc(indices, func(a, b int) int {
		return rolls[b].value - rolls[a].value // sort in descending order
	})

	for _, indx := range indices[khAction.keepCount:] {
		rolls[indx].dropped = true
	}

	return rolls
}

type DropLowest struct {
	dropCount int
}

func (postAction DropLowest) String() string {
	return fmt.Sprintf("dl%d", postAction.dropCount)
}

func (dlAction DropLowest) ApplyFilter(rolls []Roll) []Roll {
	indices := make([]int, len(rolls))
	for i := range indices {
		indices[i] = i
	}

	slices.SortFunc[[]int](indices, func(a, b int) int {
		return rolls[a].value - rolls[b].value
	})

	for _, indx := range indices[:dlAction.dropCount] {
		rolls[indx].dropped = true
	}

	return rolls
}
