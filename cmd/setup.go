package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/kumneger0/yt-tracks/internal/backend"
	"github.com/spf13/cobra"
)

func newSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Download and set up the local backend runtime",
		Long:  "Download and initialize the local Python runtime and dependencies in advance so yt-tracks launches instantly.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancel()

			if err := backend.SetupEmbedded(ctx); err != nil {
				return fmt.Errorf("error during setup: %w", err)
			}
			return nil
		},
	}
}
