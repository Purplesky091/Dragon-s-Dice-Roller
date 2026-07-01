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

func (keepHighestAction KeepHighest) ApplyFilter(rolls []int) ([]int, []int) {
	slices.Sort(rolls)
	slices.Reverse(rolls)
	kept := rolls[0:keepHighestAction.keepCount]
	dropped := rolls[keepHighestAction.keepCount:]

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
