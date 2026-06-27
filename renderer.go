package main

import (
	"fmt"
	"strconv"
	"strings"
)

func displayDiceRoll(diceStr string, diceRoll DiceRoll) string {
	sumStr := strconv.Itoa(diceRoll.result)

	if len(diceRoll.rolls) > MaxDisplayableRolls {
		return displaySingleCol(diceStr, sumStr)
	}

	return displayTwoCol(diceStr, sumStr, fmt.Sprintf("%v", diceRoll.rolls))
}

func displaySingleCol(diceStr string, sumStr string) string {
	width := max(len(diceStr), len("sum"), len(sumStr))

	top := "╔" + strings.Repeat("═", width+2) + "╗"
	mid := "╠" + strings.Repeat("═", width+2) + "╢"
	divider := "╟" + strings.Repeat("─", width+2) + "╢"
	bottom := "╚" + strings.Repeat("═", width+2) + "╝"

	return strings.Join([]string{
		"```",
		top,
		fmt.Sprintf("║ %s ║", centerPad(diceStr, width)),
		mid,
		fmt.Sprintf("║ %s ║", centerPad("sum", width)),
		divider,
		fmt.Sprintf("║ %s ║", centerPad(sumStr, width)),
		bottom,
		"```",
	}, "\n")
}

func displayTwoCol(diceStr, sumStr, rollsStr string) string {
	w1 := max(len("rolls"), len(rollsStr))
	w2 := max(len("sum"), len(sumStr))
	innerWidth := w1 + w2 + 3
	hw := max(len(diceStr), innerWidth)

	top := "╔" + strings.Repeat("═", hw+2) + "╗"
	mid := "╠" + strings.Repeat("═", w1+2) + "╤" + strings.Repeat("═", w2+2) + "╣"
	div := "╟" + strings.Repeat("─", w1+2) + "┼" + strings.Repeat("─", w2+2) + "╢"
	bot := "╚" + strings.Repeat("═", w1+2) + "╧" + strings.Repeat("═", w2+2) + "╝"

	return strings.Join([]string{
		"```",
		top,
		fmt.Sprintf("║ %s ║", centerPad(diceStr, hw)),
		mid,
		fmt.Sprintf("║ %s │ %-*s ║", centerPad("rolls", w1), w2, "sum"),
		div,
		fmt.Sprintf("║ %s │ %s ║", centerPad(rollsStr, w1), centerPad(sumStr, w2)),
		bot,
		"```",
	}, "\n")
}

func centerPad(s string, width int) string {
	pad := width - len(s)
	if pad <= 0 {
		return s
	}
	left := pad / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", pad-left)
}
