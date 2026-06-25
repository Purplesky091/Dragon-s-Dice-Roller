package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

var opts = &slog.HandlerOptions{Level: slog.LevelInfo}
var logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
var useDiscordBot = true

func main() {
	slog.SetDefault(logger)
	if !useDiscordBot {
		RunConsoleVersion()
		return
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		slog.Error("DISCORD_TOKEN env var not set.")
		os.Exit(1)
	}

}

func RunConsoleVersion() {
	scanner := bufio.NewScanner(os.Stdin)
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
	fmt.Println(diceRoll)
	winningRoll, advRolls := dice.RollAdvantage()

	fmt.Println("Winning roll:", winningRoll.result)
	for i := 0; i < len(advRolls); i++ {
		fmt.Println(advRolls[i])
	}

	fmt.Println("Rolling disadvantage")
	minRoll, disAdvRolls := dice.RollDisadvantage()
	fmt.Println("Winning roll:", minRoll.result)
	for i := 0; i < len(disAdvRolls); i++ {
		fmt.Println(disAdvRolls[i])
	}
}
