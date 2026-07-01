package main

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"regexp"
	"strconv"
)

var diceRegex = regexp.MustCompile(`^(?P<DiceCount>\d+)?d(?P<FaceCount>\d+)(?P<PostApplyFlag>kh|dl)?(?P<PostApplyCount>\d+)?$`)

const diceCountCap int = 1000
const faceCap int = 1000
const diceLengthCap int = 16

func randRange(min int, max int) int {
	return rand.IntN(max+1-min) + min
}

type Dice struct {
	count      int
	faces      int
	postAction PostAction
}

func (dice Dice) String() string {
	if dice.postAction == nil {
		return fmt.Sprintf("%dd%d", dice.count, dice.faces)
	}
	return fmt.Sprintf("%dd%d%s", dice.count, dice.faces, dice.postAction)

}

func (dice Dice) Roll() DiceRoll {
	slog.Debug("Roll", "dice", dice)

	rolls := make([]int, dice.count)
	var dropped []int = nil

	for i := 0; i < dice.count; i++ {
		roll := randRange(1, dice.faces)
		slog.Debug("Rolled value", "roll", roll)
		rolls[i] = roll
	}

	kept := rolls

	if dice.postAction != nil {
		kept, dropped = dice.postAction.ApplyFilter(rolls)
	}

	result := 0
	for _, keptRoll := range kept {
		result += keptRoll
	}

	return DiceRoll{result: result, rolls: kept, dropped: dropped}
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

	postApplyFlag := matches["PostApplyFlag"]
	var postAction PostAction = nil

	if postApplyFlag != "" {
		var postApplyCount int = 1
		slog.Info("PostApplyFlag set", "postApplyFlag", postApplyFlag)
		if matches["PostApplyCount"] != "" {
			postApplyCount, _ = strconv.Atoi(matches["PostApplyCount"])
			slog.Info("PostApplyCount set", "PostApplyCount", postApplyCount)
		} else {
			slog.Info("No PostApplyCount set. Defaulting to 1")
		}

		switch postApplyFlag {
		case "kh":
			postAction = KeepHighest{keepCount: postApplyCount}
		case "dl":
			postAction = DropLowest{dropCount: postApplyCount}
		default:
			postAction = nil
		}
	}

	return Dice{count: diceCount, faces: faceCount, postAction: postAction}, nil
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
