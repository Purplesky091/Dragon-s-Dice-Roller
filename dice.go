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

type Dice struct {
	count int
	faces int
}

func (dice Dice) String() string {
	return fmt.Sprintf("%dd%d", dice.count, dice.faces)
}

func (dice Dice) Roll() int {
	slog.Debug("Roll", "dice", dice)
	var result int = 0
	for i := 0; i < dice.count; i++ {
		roll := randRange(1, dice.faces)
		slog.Debug("Rolled value", "roll", roll)
		result += roll
	}
	return result
}

func randRange(min int, max int) int {
	return rand.IntN(max+1-min) + min
}

func (dice Dice) RollDisadvantage() (int, [2]int) {
	slog.Debug("Printing dice used", "dice", dice)
	min := math.MaxInt
	var rolls [2]int
	for i := 0; i < 2; i++ {
		rolls[i] = dice.Roll()
		slog.Info("Disadvantage roll", "dice_roll", i+1, "roll", rolls[i])
		if rolls[i] < min {
			min = rolls[i]
		}
	}

	slog.Info("Disadvantage Result", "min", min)

	return min, rolls
}

func (dice Dice) RollTripleDisadvantage() (int, [3]int) {
	slog.Debug("Printing dice used", "dice", dice)
	min := math.MaxInt
	var rolls [3]int
	for i := 0; i < 3; i++ {
		rolls[i] = dice.Roll()
		slog.Info("Disadvantage roll", "dice_roll", i+1, "roll", rolls[i])
		if rolls[i] < min {
			min = rolls[i]
		}
	}

	slog.Info("Disadvantage Result", "min", min)

	return min, rolls
}

func (dice Dice) RollAdvantage() (int, [2]int) {
	max := 0
	var rolls [2]int
	for i := 0; i < 2; i++ {
		rolls[i] = dice.Roll()
		slog.Info("Advantage roll", "dice_roll", i+1, "roll", rolls[i])
		if rolls[i] > max {
			max = rolls[i]
		}
	}

	slog.Info("Advantage Result", "max", max)

	return max, rolls
}

func (dice Dice) RollTripleAdvantage() (int, [3]int) {
	max := 0
	var rolls [3]int
	for i := 0; i < 3; i++ {
		rolls[i] = dice.Roll()
		slog.Info("Triple Advantage roll", "dice_roll", i+1, "roll", rolls[i])
		if rolls[i] > max {
			max = rolls[i]
		}
	}

	slog.Info("Triple Advantage Result", "max", max)

	return max, rolls
}

func NewDice(dice string) (Dice, error) {
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
