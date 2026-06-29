package main

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
)

var diceRegex = regexp.MustCompile(`^(?P<DiceCount>\d+)?d(?P<FaceCount>\d+)(?P<KeepFlag>kh)?(?P<KeepCount>\d+)?$`)

const diceCountCap int = 1000
const faceCap int = 1000
const diceLengthCap int = 10

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

func NewDice(dice string) (Dice, error) {
	if len(dice) > diceLengthCap {
		return Dice{}, fmt.Errorf("Dice input is too long. It must be %d characters or less", diceLengthCap)
	}

	matches, matchingErr := makeSubmatchMap(diceRegex, dice)
	if matchingErr != nil {
		return Dice{}, matchingErr
	}
	var diceCount int = 1 // assume default of 1 in case first parameter is not provided (since the number before d is optional and if left blank means 1)
	var faceCount int
	var keepCount int = -1

	if matches["DiceCount"] != "" {
		var err error
		diceCount, err = strconv.Atoi(matches["DiceCount"])
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

	faceCount, err := strconv.Atoi(matches["FaceCount"])
	if err != nil {
		return Dice{}, fmt.Errorf("Invalid face count: %w", err)
	}

	if faceCount == 0 {
		return Dice{}, fmt.Errorf("Can't roll a d0. Number of faces must be 1 or higher")
	} else if faceCount > faceCap {
		return Dice{}, fmt.Errorf("Can't roll a d%d. Max is d%d", faceCount, faceCap)
	}

	hasKeepFlag := matches["KeepFlag"] != ""
	if hasKeepFlag {
		slog.Info("Found kh")

		if matches["KeepCount"] == "" {
			keepCount = 1
			slog.Info("KeepCount set to default value")
		} else {
			keepCount, _ = strconv.Atoi(matches["KeepCount"])
			slog.Info("KeepCount set", "keepCount", keepCount)
		}
	}

	return Dice{count: diceCount, faces: faceCount}, nil
}

func makeSubmatchMap(regex *regexp.Regexp, inputString string) (map[string]string, error) {
	matches := regex.FindStringSubmatch(inputString)
	if matches == nil {
		return nil, fmt.Errorf("%q isn't a valid dice. Format is <number>d<number> (e.g. 2d6)", inputString)
	}

	subMatchMap := make(map[string]string)
	for index, captureGroup := range regex.SubexpNames() {
		if index > 0 {
			subMatchMap[captureGroup] = matches[index]
		}
	}

	return subMatchMap, nil
}
