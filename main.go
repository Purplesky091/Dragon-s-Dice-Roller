package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
)

var diceRegex = regexp.MustCompile(`^(\d+)d(\d+)$`)

type Dice struct {
	count int
	faces int
}

func (dice Dice) String() string {
	return fmt.Sprintf("{count: %d, faces: %d}", dice.count, dice.faces)
}

func (dice Dice) Roll() int {
	roll := randRange(1, dice.faces)
	fmt.Printf("Rolled a %d\n", roll)
	fmt.Printf("Multiplying %d with the count, %d\n", roll, dice.count)
	return dice.count * roll
}

func randRange(min int, max int) int {
	return rand.IntN(max+1-min) + min
}

func parseDice(dice string) (Dice, error) {
	matches := diceRegex.FindStringSubmatch(dice)
	if matches == nil {
		return Dice{}, fmt.Errorf("%q is an invalid dice. Dice must match the format of <number>d<number>. <number> must be positive too.", dice)
	}

	diceCount, err := strconv.Atoi(matches[1])
	if err != nil {
		return Dice{}, fmt.Errorf("Invalid dice count: %w", err)
	}
	faceCount, err := strconv.Atoi(matches[2])
	if err != nil {
		return Dice{}, fmt.Errorf("Invalid face count: %w", err)
	}

	return Dice{count: diceCount, faces: faceCount}, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	fmt.Print("Enter dice you wish to roll (format: <num>d<num>): ")
	if scanner.Scan() {
		input = scanner.Text()
	}
	fmt.Printf("Dice rolled: %s\n", input)
	dice, err := parseDice(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := dice.Roll()
	fmt.Printf("You rolled %d\n", result)
}
