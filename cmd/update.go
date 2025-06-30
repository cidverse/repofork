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
			originBranch, _ := cmd.Flags().GetString("origin-branch")
			upstream, _ := cmd.Flags().GetString("upstream")
			upstreamBranch, _ := cmd.Flags().GetString("upstream-branch")
			fullRewrite, _ := cmd.Flags().GetBool("full-rewrite")
			push, _ := cmd.Flags().GetBool("push")
			slog.With("origin", origin).With("upstream", upstream).Info("update fork")

			err := fork.UpdateFork(origin, originBranch, upstream, upstreamBranch, fullRewrite, push, nil)
			if err != nil {
				slog.Error("failed to update fork", "error", err)
				return
			}
		},
	}

	cmd.Flags().StringP("origin", "r", "", "origin")
	cmd.Flags().String("origin-branch", "main", "origin branch to checkout")
	cmd.Flags().StringP("upstream", "u", "", "upstream repository URL")
	cmd.Flags().String("upstream-branch", "main", "upstream branch to checkout")
	cmd.Flags().BoolP("full-rewrite", "f", false, "full rewrite of the fork")
	cmd.Flags().BoolP("push", "p", false, "push changes to origin")

	return cmd
}
