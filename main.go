package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
)

var diceRegex = regexp.MustCompile(`^(\d+)?d(\d+)$`)

type Dice struct {
	count int
	faces int
}

func (dice Dice) String() string {
	return fmt.Sprintf("{count: %d, faces: %d}", dice.count, dice.faces)
}

func (dice Dice) Roll() int {
	roll := randRange(1, dice.faces)
	return dice.count * roll
}

func randRange(min int, max int) int {
	return rand.IntN(max+1-min) + min
}

func (dice Dice) RollAdvantage() int {
	result1 := dice.Roll()
	result2 := dice.Roll()

	fmt.Println("result1:", result1)
	fmt.Println("result2:", result2)
	return max(result1, result2)
}

func (dice Dice) RollTripleAdvantage() int {
	result1 := dice.Roll()
	result2 := dice.Roll()
	result3 := dice.Roll()

	fmt.Println("result1:", result1)
	fmt.Println("result2:", result2)
	fmt.Println("result3:", result3)
	return max(result1, result2, result3)
}

func parseDice(dice string) (Dice, error) {
	matches := diceRegex.FindStringSubmatch(dice)
	var diceCount int
	var faceCount int
	if matches == nil {
		return Dice{}, fmt.Errorf("%q is an invalid dice. Dice must match the format of <number>d<number>. <number> must be positive too.", dice)
	}

	if matches[1] == "" {
		diceCount = 1
	} else {
		var err error
		diceCount, err = strconv.Atoi(matches[1])
		if err != nil {
			return Dice{}, fmt.Errorf("Invalid dice count: %w", err)
		}
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

	fmt.Println("Rolling with advantage")
	advantageResult := dice.RollAdvantage()
	fmt.Println("You rolled: ", advantageResult)

	fmt.Println("Rolling with triple advantage")
	tripleAdvantageRoll := dice.RollTripleAdvantage()
	fmt.Println("You rolled: ", tripleAdvantageRoll)
}
