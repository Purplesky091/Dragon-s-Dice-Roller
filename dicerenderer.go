package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type DiceRenderer struct {
	table   *tablewriter.Table
	builder *strings.Builder
}

func NewDiceRenderer() *DiceRenderer {
	var builder strings.Builder

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
				Alignment:    tw.CellAlignment{Global: tw.AlignCenter},
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal},
				ColMaxWidths: tw.CellWidth{Global: 17},
			},
		}),
	)
	return &DiceRenderer{table: table, builder: &builder}
}

func (diceRenderer *DiceRenderer) renderRoll(diceStr string, diceRoll DiceRoll) string {
	table := diceRenderer.table
	builder := diceRenderer.builder

	builder.Reset()
	table.Reset()
	builder.WriteString("```")

	if len(diceRoll.rolls) > MaxDisplayableRolls {
		table.Header(diceStr)
		table.Append([]string{"sum"})
		table.Append([]string{strconv.Itoa(diceRoll.result)})
	} else {
		table.Header(diceStr, diceStr)
		table.Append([]string{"rolls", "sum"})
		table.Append([]string{fmt.Sprintf("%v", diceRoll.rolls), strconv.Itoa(diceRoll.result)})
	}

	table.Render()
	builder.WriteString("```")
	return builder.String()
}
