package main

import (
	"github.com/cidverse/repofork/cmd"
	"github.com/rs/zerolog/log"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	status  = "clean"
)

// Init Hook
func init() {
	// version information
	cmd.Version = version
	cmd.CommitHash = commit
	cmd.RepositoryStatus = status
	cmd.BuildAt = date
}

// CLI Main Entrypoint
func main() {
	cmdErr := cmd.Execute()
	if cmdErr != nil {
		log.Fatal().Err(cmdErr).Msg("internal cli library error")
	}
}
