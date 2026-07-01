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

type DropLowest struct {
	dropCount int
}

func (postAction DropLowest) String() string {
	return fmt.Sprintf("dl%d", postAction.dropCount)
}

func (keepHighestAction DropLowest) ApplyFilter(rolls []int) ([]int, []int) {
	slices.Sort(rolls)
	kept := rolls[keepHighestAction.dropCount:]
	dropped := rolls[0:keepHighestAction.dropCount]
	return kept, dropped
}
