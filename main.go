package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

var opts = &slog.HandlerOptions{Level: slog.LevelInfo}
var logger = slog.New(slog.NewTextHandler(os.Stdout, opts))

func main() {
	slog.SetDefault(logger)

	scanner := bufio.NewScanner(os.Stdin)
	var input string
	fmt.Print("Enter dice you wish to roll (format: <num>d<num>): ")
	if scanner.Scan() {
		input = scanner.Text()
	}
	slog.Info("Dice rolled", "input", input)
	dice, err := NewDice(input)

	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	result := dice.Roll()
	fmt.Printf("You rolled %d\n", result)

	fmt.Println("Rolling with advantage")
	advantageResult, rolls := dice.RollAdvantage()
	fmt.Println("Advantage rolls: ", rolls)
	fmt.Println("Advantage result: ", advantageResult)

	fmt.Println("Rolling with triple advantage")
	tripleAdvantageRoll, tripleRolls := dice.RollTripleAdvantage()
	fmt.Println("Triple Advantage rolls: ", tripleRolls)
	fmt.Println("Triple Advantage result: ", tripleAdvantageRoll)

	fmt.Println("Rolling with disadvantage")
	disadvantageResult, disRolls := dice.RollDisadvantage()
	fmt.Println("Disadvantage rolls: ", disRolls)
	fmt.Println("Disadvantage result: ", disadvantageResult)

	fmt.Println("Rolling with triple disadvantage")
	trplDisResult, trpldisRolls := dice.RollTripleAdvantage()
	fmt.Println("Disadvantage rolls: ", trpldisRolls)
	fmt.Println("Disadvantage result: ", trplDisResult)
}
