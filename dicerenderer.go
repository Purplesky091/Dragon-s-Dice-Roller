package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
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

func NewDiceRenderer(RowRollSize int) *DiceRenderer {
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
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
		}),
	)

	if RowRollSize == 0 {
		RowRollSize = 5
	}

	return &DiceRenderer{table: table, builder: &builder, RowRollSize: RowRollSize}
}

func (diceRenderer *DiceRenderer) createRollsSubtable(rolls []Roll) string {
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
			if roll.dropped {
				rowRoll[j] = fmt.Sprintf("(%s)", roll)
			} else {
				rowRoll[j] = roll.String()
			}
		}

		table.Append(rowRoll)
	}

	table.Render()
	return buf.String()
}

func (diceRenderer *DiceRenderer) RenderEmbed(diceStr string, rollResult RollResult) *discordgo.MessageEmbed {
	color := 0x5865F2

	if len(rollResult.rolls) > MaxDisplayableRolls {
		return &discordgo.MessageEmbed{
			Title: diceStr,
			Color: color,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Total", Value: strconv.Itoa(rollResult.sum)},
			},
		}
	}

	var rollsBuilder strings.Builder
	for i, roll := range rollResult.rolls {
		if i > 0 {
			if i%diceRenderer.RowRollSize == 0 {
				rollsBuilder.WriteString("\n")
			} else {
				rollsBuilder.WriteString("  ")
			}
		}
		if roll.dropped {
			fmt.Fprintf(&rollsBuilder, "~~%d~~", roll.value)
		} else {
			rollsBuilder.WriteString(strconv.Itoa(roll.value))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: diceStr,
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Rolls", Value: rollsBuilder.String(), Inline: true},
			{Name: "Total", Value: strconv.Itoa(rollResult.sum), Inline: true},
		},
	}

	if rollResult.hasDroppedValues {
		embed.Footer = &discordgo.MessageEmbedFooter{Text: "Crossed out values were dropped"}
	}

	return embed
}

func (diceRenderer *DiceRenderer) RenderRoll(diceStr string, rollResult RollResult) string {
	table := diceRenderer.table
	builder := diceRenderer.builder

	builder.Reset()
	table.Reset()
	builder.WriteString("```\n")

	if len(rollResult.rolls) > MaxDisplayableRolls {
		table.Header(diceStr)
		table.Append([]string{"sum"})
		table.Append([]string{strconv.Itoa(rollResult.sum)})
	} else {
		table.Header(diceStr, diceStr)
		table.Append([]string{"rolls", "sum"})
		table.Append([]string{diceRenderer.createRollsSubtable(rollResult.rolls), strconv.Itoa(rollResult.sum)})
	}

	table.Render()
	builder.WriteString("```")
	builder.WriteString("Dropped values surrounded in ()\n")
	return builder.String()
}
