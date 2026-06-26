package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var opts = &slog.HandlerOptions{Level: slog.LevelInfo}
var logger = slog.New(slog.NewTextHandler(os.Stdout, opts))

const useDiscordBot = true
const MaxDiscordMsgLength = 2000
const MaxDisplayableRolls = 200 // how big the dice count can be before I stop showing the rolls

var rollOptions = []*discordgo.ApplicationCommandOption{
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "dice",
		Description: "Dice format (i.e. 2d6, d20, 4d4)",
		Required:    true,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "type",
		Description: "What type of roll (advantage, disadvantage, normal). Defaults to normal",
		Required:    false,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{Name: "normal", Value: "normal"},
			{Name: "advantage", Value: "advantage"},
			{Name: "disadvantage", Value: "disadvantage"},
		},
	},
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "roll",
		Description: "Roll dice in <number>d<number> format (ie 2d6)",
		Options:     rollOptions,
	},
	{
		Name:        "dr",
		Description: "alias for dice roll",
		Options:     rollOptions,
	},
}

func respond(session *discordgo.Session, interactionEvent *discordgo.InteractionCreate, content string) {
	err := session.InteractionRespond(interactionEvent.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		slog.Error("Failed to respond to interaction", "error", err)
	}
}

// oooh amazing comment

func createDiceRollMessage(dice Dice, diceRoll DiceRoll) string {
	rollsStr := ""
	if dice.count <= MaxDisplayableRolls {
		rollsStr = fmt.Sprintf("Rolls: %v\n", diceRoll.rolls)
	}
	return fmt.Sprintf("Rolling %s\n%sSum: %d", dice, rollsStr, diceRoll.result)
}

func handleRoll(session *discordgo.Session, interactionEvent *discordgo.InteractionCreate) {
	options := interactionEvent.ApplicationCommandData().Options

	diceStr := options[0].StringValue()
	dice, err := NewDice(diceStr)
	if err != nil {
		respond(session, interactionEvent, fmt.Sprintf("Invalid dice: %s", err))
		return
	}
	rollType := "normal"
	if len(options) > 1 {
		rollType = options[1].StringValue()
	}

	var msg string
	switch rollType {
	case "advantage":
		result, rolls := dice.RollAdvantage()
		msg = fmt.Sprintf("Rolling %s with advantage\nRoll 1:\n%s\nRoll 2:\n%s\nResulting Roll:\n%s", dice, rolls[0].DiscordString(), rolls[1].DiscordString(), result.DiscordString())
	case "disadvantage":
		result, rolls := dice.RollDisadvantage()
		msg = fmt.Sprintf("Rolling %s with disadvantage\nRoll 1:\n%s\nRoll 2:\n%s\nResulting Roll:\n%s", dice, rolls[0].DiscordString(), rolls[1].DiscordString(), result.DiscordString())
	default:
		roll := dice.Roll()
		msg = fmt.Sprintf("Rolling %s\n%s", dice, roll.DiscordString())
	}

	respond(session, interactionEvent, msg)
}

func main() {
	slog.SetDefault(logger)
	if !useDiscordBot {
		cli := new(CLI)
		cli.Run()
		return
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		slog.Error("DISCORD_TOKEN env var not set.")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord session", "error", err)
		os.Exit(1)
	}

	dg.AddHandler(func(session *discordgo.Session, interactionEvent *discordgo.InteractionCreate) {
		if interactionEvent.Type != discordgo.InteractionApplicationCommand {
			return
		}

		switch interactionEvent.ApplicationCommandData().Name {
		case "roll", "dr":
			handleRoll(session, interactionEvent)
		}
	})

	err = dg.Open()
	if err != nil {
		slog.Error("Error opening connection to Discord", "error", err)
		os.Exit(1)
	}
	defer dg.Close()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for index, cmd := range commands {
		registered, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd)
		if err != nil {
			slog.Error("Could not register command", "name", cmd.Name)
			os.Exit(1)
		}
		registeredCommands[index] = registered
		slog.Info("Registered slash command", "name", cmd.Name)
	}

	slog.Info("Bot is running. Press ctrl + c to quit")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// Clean up slash commands on shutdown so they don't linger in Discord.
	for _, cmd := range registeredCommands {
		if err := dg.ApplicationCommandDelete(dg.State.User.ID, "", cmd.ID); err != nil {
			slog.Warn("Could not delete command on shutdown", "name", cmd.Name, "error", err)
		}
	}
	slog.Info("Bot shut down cleanly.")
}
