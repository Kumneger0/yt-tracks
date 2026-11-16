package main

import (
	"log/slog"
	"os"

	"github.com/kumneger0/yt-tracks/cmd"
)

var (
	version = ""
	Debug   = "false"
)

func main() {
	err := cmd.Execute(version, Debug == "true")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
