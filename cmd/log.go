package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kumneger0/yt-tracks/internal/config"
	"github.com/spf13/cobra"
)

func ytTracksLog() *cobra.Command {
	userConfig := config.GetUserConfig(runtime.GOOS)
	return &cobra.Command{

		Use:          "log",
		Short:        "show log files use it for error reporting",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			if userConfig == nil {
				fmt.Println("No user config found.")
				return
			}
			if userConfig.DebugDir == nil {
				fmt.Println("No debug directory configured.")
				return
			}
			logs, err := os.ReadFile(filepath.Join(*userConfig.DebugDir, "yt-tracks.log"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading log file: %v", err)
				os.Exit(1)
			}
			fmt.Print(string(logs))
		},
	}
}
