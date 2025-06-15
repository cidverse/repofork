package cmd

import (
	"log/slog"

	"github.com/cidverse/repofork/pkg/fork"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{},
		Short:   `update a fork`,
		Run: func(cmd *cobra.Command, args []string) {
			origin, _ := cmd.Flags().GetString("origin")
			upstream, _ := cmd.Flags().GetString("upstream")
			slog.With("origin", origin).With("upstream", upstream).Info("update fork")

			err := fork.UpdateFork(origin, upstream)
			if err != nil {
				slog.Error("failed to update fork", "error", err)
				return
			}
		},
	}

	cmd.Flags().StringP("origin", "r", "", "origin")
	cmd.Flags().StringP("upstream", "u", "", "upstream repository URL")

	return cmd
}
