package main

import (
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type DiceRenderer struct {
	table       *tablewriter.Table
	builder     *strings.Builder
	RowRollSize int
}

// const RowRollSize int = 5

func NewDiceRenderer(RowRollSize ...int) *DiceRenderer {
	var builder strings.Builder
	var finalRowRollSize int = 5

	table := tablewriter.NewTable(&builder, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
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

	if len(RowRollSize) > 0 {
		finalRowRollSize = RowRollSize[0]
	}

	return &DiceRenderer{table: table, builder: &builder, RowRollSize: finalRowRollSize}
}

func (diceRenderer *DiceRenderer) createRollsSubtable(rolls []int) string {
	var buf strings.Builder
	table := tablewriter.NewTable(&buf,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.SeparatorsNone,
				Lines:      tw.LinesNone,
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}},
		}),
	)

	for i := 0; i < len(rolls); i += diceRenderer.RowRollSize {
		end := min(i+diceRenderer.RowRollSize, len(rolls))
		rowRoll := make([]string, end-i)
		for j, roll := range rolls[i:end] {
			rowRoll[j] = strconv.Itoa(roll)
		}

		table.Append(rowRoll)
	}

	table.Render()
	return buf.String()
}

func (diceRenderer *DiceRenderer) RenderRoll(diceStr string, diceRoll DiceRoll) string {
	table := diceRenderer.table
	builder := diceRenderer.builder

	builder.Reset()
	table.Reset()
	builder.WriteString("```\n")

	if len(diceRoll.rolls) > MaxDisplayableRolls {
		table.Header(diceStr)
		table.Append([]string{"sum"})
		table.Append([]string{strconv.Itoa(diceRoll.result)})
	} else if len(diceRoll.dropped) > 0 {
		table.Header(diceStr, diceStr, diceStr)
		table.Append([]string{"rolls", "dropped", "sum"})
		table.Append([]string{diceRenderer.createRollsSubtable(diceRoll.rolls), diceRenderer.createRollsSubtable(diceRoll.dropped), strconv.Itoa(diceRoll.result)})
	} else {
		table.Header(diceStr, diceStr)
		table.Append([]string{"rolls", "sum"})
		table.Append([]string{diceRenderer.createRollsSubtable(diceRoll.rolls), strconv.Itoa(diceRoll.result)})
	}

	table.Render()
	builder.WriteString("```")
	return builder.String()
}
