package main

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
)

var diceRegex = regexp.MustCompile(`^(\d+)?d(\d+)$`)
var diceCountCap = 1000
var faceCap = 1000
var diceLengthCap = 10

func randRange(min int, max int) int {
	return rand.IntN(max+1-min) + min
}

type Dice struct {
	count int
	faces int
}

func (dice Dice) String() string {
	return fmt.Sprintf("%dd%d", dice.count, dice.faces)
}

func (dice Dice) Roll() DiceRoll {
	slog.Debug("Roll", "dice", dice)
	var result int = 0
	rolls := make([]int, dice.count)

	for i := 0; i < dice.count; i++ {
		roll := randRange(1, dice.faces)
		slog.Debug("Rolled value", "roll", roll)
		rolls[i] = roll
		result += roll
	}
	return DiceRoll{result: result, rolls: rolls}
}

func (dice Dice) RollAdvantage() (DiceRoll, []DiceRoll) {
	var maxDice DiceRoll = DiceRoll{result: 0, rolls: nil}
	var diceRolls []DiceRoll = make([]DiceRoll, 2)

	for i := 0; i < 2; i++ {
		diceRolls[i] = dice.Roll()
		if diceRolls[i].result > maxDice.result {
			maxDice = diceRolls[i]
		}
	}

	slog.Info("Advantage roll result", "maxDice", maxDice)

	return maxDice, diceRolls
}

func (dice Dice) RollDisadvantage() (DiceRoll, []DiceRoll) {
	slog.Debug("Printing dice used", "dice", dice)
	var minRoll DiceRoll = DiceRoll{result: math.MaxInt, rolls: nil}
	var rolls []DiceRoll = make([]DiceRoll, 2)

	for i := 0; i < 2; i++ {
		rolls[i] = dice.Roll()
		slog.Info("Disadvantage roll", "dice_roll", i+1, "roll", rolls[i])
		if rolls[i].result < minRoll.result {
			minRoll = rolls[i]
		}
	}

	slog.Info("Disadvantage roll result", "minRoll", minRoll)

	return minRoll, rolls
}

// func (dice Dice) RollTripleDisadvantage() (int, [3]int) {
// 	slog.Debug("Printing dice used", "dice", dice)
// 	min := math.MaxInt
// 	var rolls [3]int
// 	for i := 0; i < 3; i++ {
// 		rolls[i] = dice.Roll()
// 		slog.Info("Disadvantage roll", "dice_roll", i+1, "roll", rolls[i])
// 		if rolls[i] < min {
// 			min = rolls[i]
// 		}
// 	}

// 	slog.Info("Disadvantage Result", "min", min)

// 	return min, rolls
// }

// func (dice Dice) RollTripleAdvantage() (int, [3]int) {
// 	max := 0
// 	var rolls [3]int
// 	for i := 0; i < 3; i++ {
// 		rolls[i] = dice.Roll()
// 		slog.Info("Triple Advantage roll", "dice_roll", i+1, "roll", rolls[i])
// 		if rolls[i] > max {
// 			max = rolls[i]
// 		}
// 	}

// 	slog.Info("Triple Advantage Result", "max", max)

// 	return max, rolls
// }

func NewDice(dice string) (Dice, error) {
	if len(dice) > diceLengthCap {
		return Dice{}, fmt.Errorf("Dice input is too long. It must be %d characters or less", diceLengthCap)
	}

	matches := diceRegex.FindStringSubmatch(dice)
	var diceCount int = 1 // assume default of 1 in case first parameter is not provided (since the number before d is optional and if left blank means 1)
	var faceCount int
	if matches == nil {
		return Dice{}, fmt.Errorf("%q isn't a valid dice. Format is <number>d<number> (e.g. 2d6)", dice)
	}

	if matches[1] != "" {
		var err error
		diceCount, err = strconv.Atoi(matches[1])
		if err != nil {
			return Dice{}, fmt.Errorf("Invalid dice count: %w", err)
		}

		if diceCount == 0 {
			return Dice{}, fmt.Errorf("Number of dice must be greater than 0")
		}

		if diceCount > diceCountCap {
			return Dice{}, fmt.Errorf("The dice count %d in %q is too high. Max is %d", diceCount, dice, diceCountCap)
		}
	}

	faceCount, err := strconv.Atoi(matches[2])
	if err != nil {
		return Dice{}, fmt.Errorf("Invalid face count: %w", err)
	}

	if faceCount == 0 {
		return Dice{}, fmt.Errorf("Can't roll a d0. Number of faces must be 1 or higher")
	} else if faceCount > faceCap {
		return Dice{}, fmt.Errorf("Can't roll a d%d. Max is d%d", faceCount, faceCap)
	}

	return Dice{count: diceCount, faces: faceCount}, nil
}
