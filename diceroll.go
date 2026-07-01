package main

import (
	"fmt"
	"strconv"
	"strings"
)

type DiceRoll struct {
	result  int
	rolls   []int
	dropped []int
}

func (diceRoll DiceRoll) String() string {
	return fmt.Sprintf("rolls:%v\nsum:%d", diceRoll.rolls, diceRoll.result)
}

func (diceRoll DiceRoll) writeRollsWrapped() string {
	var builder strings.Builder
	for index, roll := range diceRoll.rolls {
		builder.WriteString(strconv.Itoa(roll))
		if (index+1)%5 == 0 {
			builder.WriteString("\n")

		} else {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}
