package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

type CLI struct{}

func (console *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	var diceRenderer *DiceRenderer = NewDiceRenderer(5)
	var input string
	fmt.Print("Enter dice you wish to roll (format: <num>d<num>): ")
	if scanner.Scan() {
		input = scanner.Text()
	}
	slog.Debug("User input", "input", input)
	dice, err := NewDice(input)

	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	diceRoll := dice.Roll()
	result := diceRenderer.RenderRoll(dice.String(), diceRoll)
	fmt.Println(result)
}
