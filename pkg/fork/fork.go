package fork

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
)

func UpdateFork(remote string, upstream string) error {
	// create a temporary directory
	tempDir, err := os.MkdirTemp("", "repofork-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	slog.Info("Created temporary directory", "dir", tempDir)

	// git init
	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// add remotes
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remote},
	})
	if err != nil {
		return fmt.Errorf("failed to create remote 'origin': %w", err)
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{upstream},
	})
	if err != nil {
		return fmt.Errorf("failed to create remote 'upstream': %w", err)
	}

	// fetch
	_ = NewGit(tempDir)
	err = gitCommand(tempDir, "fetch", "origin")
	if err != nil {
		return err
	}
	err = gitCommand(tempDir, "fetch", "upstream")
	if err != nil {
		return err
	}

	// checkout main branch
	if err = gitCommand(tempDir, "checkout", "-B", "main", "origin/main"); err != nil {
		return err
	}
	if err = gitCommand(tempDir, "branch", "--track", "upstream", "upstream/main"); err != nil {
		return err
	}

	// get last mirrored upstream commit SHA
	lastSHA, err := getLastUpstreamCommit(tempDir)
	if err != nil {
		return fmt.Errorf("failed to get last upstream commit: %w", err)
	}
	slog.Info("Last upstream commit SHA", "sha", lastSHA)

	// cherry-pick
	if lastSHA == "" {
		slog.Info("No previous mirrored commit found — full rewrite")

		if err = gitCommand(tempDir, "checkout", "-B", "main", "upstream/main"); err != nil {
			return fmt.Errorf("failed to checkout upstream/main: %w", err)
		}
	} else {
		slog.Info("Last mirrored upstream commit found", "sha", lastSHA)

		// Create temp branch from last known upstream commit
		if err = gitCommand(tempDir, "checkout", "-B", "main", "origin/main"); err != nil {
			return fmt.Errorf("failed to checkout origin/main: %w", err)
		}

		// Cherry-pick new commits from upstream/main
		if err = gitCommand(tempDir, "cherry-pick", lastSHA+"..upstream/main"); err != nil {
			return fmt.Errorf("failed to cherry-pick upstream changes: %w", err)
		}
	}

	// filter repo history
	err = gitCommand(tempDir, "filter-repo",
		"--commit-callback", `if b"Original-Upstream-Commit:" in commit.message:
	commit.skip()
else:
	commit.message += b"\nOriginal-Upstream-Commit: " + commit.original_id`,
		"--invert-paths",
		"--path", ".github/workflows/",
		"--path", ".gitlab-ci.yml",
		"--path", ".github/",
		"--path", ".gitlab/",
		"--force",
	)

	// open after history rewrite
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remote},
	})
	if err != nil {
		return fmt.Errorf("failed to create remote 'origin': %w", err)
	}

	// push
	//if err = gitCommand(tempDir, "push", "--force", "origin", "main"); err != nil {
	//	return fmt.Errorf("failed to push changes to origin: %w", err)
	//}

	return nil
}

func getLastUpstreamCommit(repoDir string) (string, error) {
	output, err := gitCommandOutput(repoDir, "log", "origin/main", "--grep=Original-Upstream-Commit:", "-n", "1")
	if err != nil {
		return "", fmt.Errorf("failed to get git log: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "    Original-Upstream-Commit:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "    Original-Upstream-Commit:")), nil
		}
	}
	return "", nil
}
