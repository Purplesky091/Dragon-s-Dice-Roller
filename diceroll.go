package main

import "fmt"

type DiceRoll struct {
	result int
	rolls  []int
}

func (diceRoll DiceRoll) String() string {
	return fmt.Sprintf("rolls:%v\nsum:%d", diceRoll.rolls, diceRoll.result)
}

func (diceRoll DiceRoll) DiscordString() string {
	result := ""
	if len(diceRoll.rolls) <= MaxDisplayableRolls {
		result = fmt.Sprintf("rolls: %v\nsum: %d", diceRoll.rolls, diceRoll.result)
	} else {
		result = fmt.Sprintf("sum: %d", diceRoll.result)
	}

	return result
}
