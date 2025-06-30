package cmd

import (
	"log/slog"
	"os"

	"github.com/cidverse/repofork/pkg/config"
	"github.com/cidverse/repofork/pkg/fork"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{},
		Short:   `update one or multiple forks based on the configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			configFile, _ := cmd.Flags().GetString("config")
			forkName, _ := cmd.Flags().GetString("name")

			conf, err := config.LoadConfig(configFile)
			if err != nil {
				slog.With("err", err).Error("error loading config")
				os.Exit(1)
			}

			if len(conf.Forks) == 0 {
				slog.Error("no forks configured")
				os.Exit(1)
			}

			for _, f := range conf.Forks {
				if forkName != "" && f.Name != forkName {
					slog.With("fork", f.Name).Debug("Skipping fork due to name mismatch")
					continue
				}
				slog.Info("Processing fork", "name", f.Name, "mirror", f.OriginRepo, "upstream", f.UpstreamRepo)

				err = fork.UpdateFork(f.OriginRepo, f.OriginBranch, f.UpstreamRepo, f.UpstreamBranch, f.FullRewrite, f.Push, f.ExcludePaths)
				if err != nil {
					slog.With("fork", f.Name).With("err", err).Error("Failed to update fork")
				} else {
					slog.Info("Successfully updated fork", "name", f.Name)
				}
			}
		},
	}

	cmd.Flags().StringP("config", "c", "", "config file")
	cmd.Flags().StringP("name", "n", "", "name of the fork to update (optional, updates all if not set)")

	return cmd
}
