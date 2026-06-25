package main

import (
	"log/slog"
	"os"
)

var opts = &slog.HandlerOptions{Level: slog.LevelInfo}
var logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
var useDiscordBot = false

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
}
