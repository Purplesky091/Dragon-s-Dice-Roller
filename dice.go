package main

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"regexp"
	"strconv"
)

const DICE_COUNT = "DiceCount"
const FACE_COUNT = "FaceCount"
const POST_ACTION_FLAG = "PostActionFlag"
const POST_ACTION_COUNT = "PostActionCount"

var diceRegex = regexp.MustCompile(
	fmt.Sprintf(`^(?P<%s>\d+)?d(?P<%s>\d+)(?P<%s>kh|kl|dh|dl)?(?P<%s>\d+)?$`,
		DICE_COUNT,
		FACE_COUNT,
		POST_ACTION_FLAG,
		POST_ACTION_COUNT))

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

	diceCount, diceCountErr := getDiceCount(matches[DICE_COUNT], dice)
	if diceCountErr != nil {
		return Dice{}, diceCountErr
	}

	faceCount, faceCountErr := getFaceCount(matches[FACE_COUNT])
	if faceCountErr != nil {
		return Dice{}, faceCountErr
	}

	postAction := getPostAction(matches[POST_ACTION_FLAG], matches[POST_ACTION_COUNT])

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

func getDiceCount(diceCountStr string, diceStr string) (int, error) {
	// assume default of 1 in case first parameter is not provided (since the number before d is optional and if left blank means 1)
	if diceCountStr == "" {
		return 1, nil
	}

	var err error
	var diceCount int

	diceCount, err = strconv.Atoi(diceCountStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid dice count: %w", err)
	}

	if diceCount == 0 {
		return 0, fmt.Errorf("Number of dice must be greater than 0")
	}

	if diceCount > diceCountCap {
		return 0, fmt.Errorf("The dice count %d in %q is too high. Max is %d", diceCount, diceStr, diceCountCap)
	}

	return diceCount, nil
}

func getFaceCount(faceCountStr string) (int, error) {
	faceCount, err := strconv.Atoi(faceCountStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid face count: %w", err)
	}

	if faceCount == 0 {
		return 0, fmt.Errorf("Can't roll a d0. Number of faces must be 1 or higher")
	} else if faceCount > faceCap {
		return 0, fmt.Errorf("Can't roll a d%d. Max is d%d", faceCount, faceCap)
	}

	return faceCount, nil
}

func getPostAction(postActionFlag string, postActionCountStr string) PostAction {
	var postAction PostAction = nil
	var postActionCount int

	if postActionFlag == "" {
		return nil
	}

	slog.Info("postActionFlag set", "postActionFlag", postActionFlag)
	if postActionCountStr != "" {
		postActionCount, _ = strconv.Atoi(postActionCountStr)
		slog.Info("postActionCount set", "postActionCount", postActionCount)
	} else {
		postActionCount = 1
		slog.Info("No postActionCount set. Defaulting to 1")
	}

	switch postActionFlag {
	case "kh":
		postAction = KeepHighest{keepCount: postActionCount}
	case "kl":
		postAction = KeepLowest{keepCount: postActionCount}
	case "dh":
		postAction = DropHighest{dropCount: postActionCount}
	case "dl":
		postAction = DropLowest{dropCount: postActionCount}

	}

	return postAction
}
