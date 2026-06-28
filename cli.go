package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type CLI struct{}

func (console *CLI) Run() {
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

	// var builder strings.Builder
	table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
		Settings: tw.Settings{
			Separators: tw.Separators{
				BetweenRows: tw.On,
			},
		},
	})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					MergeMode:  tw.MergeHorizontal,
					AutoFormat: tw.Off,
				},
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
		}),
	)

	table.Reset()
	table.Header(dice.String(), dice.String())
	table.Append([]string{"rolls", "sum"})
	table.Append([]string{fmt.Sprintf("%v", diceRoll.rolls), strconv.Itoa(diceRoll.result)})
	table.Render()

	table.Reset()
	table.Header("d20", "d20")
	table.Append([]string{"rolls", "sum"})
	table.Append([]string{"[10]", "10"})
	table.Render()
}
