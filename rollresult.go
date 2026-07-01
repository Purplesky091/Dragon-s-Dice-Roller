package main

import (
	"fmt"
)

type RollResult struct {
	sum              int
	rolls            []Roll
	hasDroppedValues bool
}

func (diceRoll RollResult) String() string {
	return fmt.Sprintf("rolls:%v\nsum:%d", diceRoll.rolls, diceRoll.sum)
}
