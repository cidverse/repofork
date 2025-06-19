package fork

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
)

func UpdateFork(remote string, originBranch string, upstream string, upstreamBranch string, fullRewrite bool, push bool) error {
	originRef := "origin/" + originBranch
	upstreamRef := "upstream/" + upstreamBranch

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
	g := NewGit(tempDir)
	err = gitCommand(tempDir, "fetch", "origin")
	if err != nil {
		return err
	}
	err = gitCommand(tempDir, "fetch", "upstream")
	if err != nil {
		return err
	}

	// checkout main branch
	lastSHA := ""
	if g.RefExists(originRef) {
		if err = gitCommand(tempDir, "checkout", "-B", "main", originRef); err != nil {
			return err
		}

		// get last mirrored upstream commit SHA
		lastSHA, err = g.LastUpstreamCommit(originRef)
		if err != nil {
			return fmt.Errorf("failed to get last upstream commit: %w", err)
		}
		slog.Info("Last upstream commit SHA", "sha", lastSHA)
	}
	if err = gitCommand(tempDir, "branch", "--track", "upstream", upstreamRef); err != nil {
		return err
	}

	// cherry-pick
	if lastSHA == "" || fullRewrite {
		slog.Info("No previous mirrored commit found — full rewrite")

		if err = gitCommand(tempDir, "checkout", "-B", "main", upstreamRef); err != nil {
			return fmt.Errorf("failed to checkout upstream/main: %w", err)
		}
	} else {
		slog.Info("Last mirrored upstream commit found", "sha", lastSHA)

		// Create temp branch from last known upstream commit
		if err = gitCommand(tempDir, "checkout", "-B", "main", originRef); err != nil {
			return fmt.Errorf("failed to checkout origin/main: %w", err)
		}

		// cherry-pick commits
		commits, err := g.CommitsBetween(lastSHA, upstreamRef)
		if err != nil {
			return fmt.Errorf("failed to get commits between %s and %s: %w", lastSHA, upstreamRef, err)
		}
		for _, sha := range commits {
			slog.Info("Commit found", "sha", sha)

			isMerge, _ := g.IsMergeCommit(sha)
			if isMerge {
				if err = gitCommand(tempDir, "cherry-pick", sha, "-m", "1"); err != nil { // parent 1 is the checked out branch when the merge was created, required for cherry-picking merge commits
					return fmt.Errorf("failed to cherry-pick commit %s: %w", sha, err)
				}
			} else {
				if err = gitCommand(tempDir, "cherry-pick", sha); err != nil {
					return fmt.Errorf("failed to cherry-pick commit %s: %w", sha, err)
				}
			}
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
	if push {
		if err = gitCommand(tempDir, "push", "--force", "origin", "HEAD:"+originBranch); err != nil {
			return fmt.Errorf("failed to push changes to origin:%s: %w", originBranch, err)
		}
	}

	return nil
}
