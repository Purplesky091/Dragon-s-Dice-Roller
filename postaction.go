package main

import (
	"fmt"
	"slices"
)

type PostAction interface {
	ApplyFilter(rolls []int) ([]int, []int)
}

type KeepHighest struct {
	keepCount int
}

func (postAction KeepHighest) String() string {
	return fmt.Sprintf("kh%d", postAction.keepCount)
}

func (khAction KeepHighest) ApplyFilter(rolls []int) ([]int, []int) {
	indices := make([]int, len(rolls))
	for i := range indices {
		indices[i] = i
	}

	// Sort a slice of indices by their corresponding roll values (descending).
	// This lets us rank rolls by value without reordering the original slice,
	// so we can mark which positions are dropped while preserving insertion order.
	slices.SortFunc(indices, func(a, b int) int {
		return rolls[b] - rolls[a] // sort in descending order
	})

	isDropped := make([]bool, len(rolls))
	for _, idx := range indices[khAction.keepCount:] {
		isDropped[idx] = true
	}

	kept := make([]int, 0, khAction.keepCount)
	dropped := make([]int, 0, len(rolls)-khAction.keepCount)
	for i, roll := range rolls {
		if isDropped[i] {
			dropped = append(dropped, roll)
		} else {
			kept = append(kept, roll)
		}
	}

	return kept, dropped
}

type KeepLowest struct {
	keepCount int
}

func (postAction KeepLowest) String() string {
	return fmt.Sprintf("kl%d", postAction.keepCount)
}

func (klAction KeepLowest) ApplyFilter(rolls []int) ([]int, []int) {
	indices := make([]int, len(rolls))
	for i := range indices {
		indices[i] = i
	}

	// Sort a slice of indices by their corresponding roll values (descending).
	// This lets us rank rolls by value without reordering the original slice,
	// so we can mark which positions are dropped while preserving insertion order.
	slices.SortFunc(indices, func(a, b int) int {
		return rolls[a] - rolls[b] // sort in descending order
	})

	isDropped := make([]bool, len(rolls))
	for _, idx := range indices[klAction.keepCount:] {
		isDropped[idx] = true
	}

	kept := make([]int, 0, klAction.keepCount)
	dropped := make([]int, 0, len(rolls)-klAction.keepCount)
	for i, roll := range rolls {
		if isDropped[i] {
			dropped = append(dropped, roll)
		} else {
			kept = append(kept, roll)
		}
	}

	return kept, dropped
}

type DropLowest struct {
	dropCount int
}

func (postAction DropLowest) String() string {
	return fmt.Sprintf("dl%d", postAction.dropCount)
}

func (dlAction DropLowest) ApplyFilter(rolls []int) ([]int, []int) {
	indices := make([]int, len(rolls))
	for i := range indices {
		indices[i] = i
	}

	slices.SortFunc[[]int](indices, func(a, b int) int {
		return rolls[a] - rolls[b]
	})

	isDropped := make([]bool, len(rolls))
	for _, indx := range indices[:dlAction.dropCount] {
		isDropped[indx] = true
	}

	kept := make([]int, 0, len(rolls)-dlAction.dropCount)
	dropped := make([]int, 0, dlAction.dropCount)

	for i, roll := range rolls {
		if isDropped[i] {
			dropped = append(dropped, roll)
		} else {
			kept = append(kept, roll)
		}
	}

	return kept, dropped
}
