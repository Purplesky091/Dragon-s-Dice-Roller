package main

import (
	"fmt"
	"slices"
)

type PostAction interface {
	ApplyFilter(rolls []int) []int
}

type KeepHighest struct {
	keepCount int
}

func (postAction KeepHighest) String() string {
	return fmt.Sprintf("kh%d", postAction.keepCount)
}

func (keepHighestAction KeepHighest) ApplyFilter(rolls []int) []int {
	slices.Sort(rolls)
	slices.Reverse(rolls)
	return rolls[0:keepHighestAction.keepCount]
}

type DropLowest struct {
	dropCount int
}

func (postAction DropLowest) String() string {
	return fmt.Sprintf("dl%d", postAction.dropCount)
}

func (keepHighestAction DropLowest) ApplyFilter(rolls []int) []int {
	slices.Sort(rolls)
	return rolls[keepHighestAction.dropCount:]
}
