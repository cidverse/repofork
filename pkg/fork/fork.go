package fork

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
)

const tempBranchName = "repofork-temp"
const workingBranchName = "main"

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
	g := NewGit(tempDir)

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

	// fetch remotes
	if err = gitCommand(tempDir, "fetch", "--all"); err != nil {
		return fmt.Errorf("failed to fetch remotes: %w", err)
	}

	// checkout main branch
	lastSHA := ""
	if g.RefExists("refs/remotes/" + originRef) {
		if err = gitCommand(tempDir, "checkout", "-B", tempBranchName, originRef); err != nil {
			return fmt.Errorf("failed to checkout origin/%s: %w", originBranch, err)
		}

		// get last mirrored upstream commit SHA
		lastSHA, err = g.LastUpstreamCommit(originRef)
		if err != nil {
			return fmt.Errorf("failed to get last upstream commit: %w", err)
		}
	}

	// cherry-pick
	if lastSHA == "" || fullRewrite {
		slog.Info("No previous mirrored commit found — full rewrite")

		if err = gitCommand(tempDir, "checkout", "-B", workingBranchName, upstreamRef); err != nil {
			return fmt.Errorf("failed to checkout upstream/main: %w", err)
		}
	} else {
		slog.Info("Last mirrored upstream commit found", "sha", lastSHA)

		if err = gitCommand(tempDir, "checkout", "-B", workingBranchName, originRef); err != nil {
			return fmt.Errorf("failed to checkout origin/%s: %w", originBranch, err)
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
				if err = gitCommand(tempDir, "cherry-pick", sha, "-m", "1", "--allow-empty"); err != nil { // parent 1 is the checked out branch when the merge was created, required for cherry-picking merge commits
					return fmt.Errorf("failed to cherry-pick commit %s: %w", sha, err)
				}
			} else {
				if err = gitCommand(tempDir, "cherry-pick", sha, "--allow-empty"); err != nil {
					return fmt.Errorf("failed to cherry-pick commit %s: %w", sha, err)
				}
			}
		}
	}

	// filter repo history
	err = gitCommand(tempDir, "filter-repo",
		"--refs", workingBranchName,
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

	// push
	if push {
		if err = gitCommand(tempDir, "push", "--force", "origin", "HEAD:"+originBranch); err != nil {
			return fmt.Errorf("failed to push changes to origin:%s: %w", originBranch, err)
		}
	}

	return nil
}
