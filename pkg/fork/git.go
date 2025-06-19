package fork

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	repoDir string
}

func (g *Git) RefExists(ref string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", ref)
	cmd.Dir = g.repoDir
	return cmd.Run() == nil
}

func (g *Git) IsMergeCommit(sha string) (bool, error) {
	output, err := gitCommandOutput(g.repoDir, "rev-list", "--parents", "-n", "1", sha)
	if err != nil {
		return false, fmt.Errorf("failed to check if commit %s is a merge: %w", sha, err)
	}

	// A merge commit has more than one parent (i.e., more than two fields: SHA + N parents)
	fields := strings.Fields(string(output))
	return len(fields) > 2, nil
}

func (g *Git) LastUpstreamCommit(ref string) (string, error) {
	output, err := gitCommandOutput(g.repoDir, "log", ref, "--grep=Original-Upstream-Commit:", "-n", "1")
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

func (g *Git) CommitsBetween(fromRef, toRef string) ([]string, error) {
	// Format: git log --reverse --pretty=format:%H fromSHA..toRef
	output, err := gitCommandOutput(
		g.repoDir,
		"log",
		"--reverse",                           // chronological order
		"--pretty=format:%H",                  // only commit hashes
		fmt.Sprintf("%s..%s", fromRef, toRef), // range
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit list: %w", err)
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(commits) == 1 && commits[0] == "" {
		return []string{}, nil // no commits
	}

	return commits, nil
}

func NewGit(repoDir string) *Git {
	return &Git{repoDir: repoDir}
}

func gitCommand(repoDir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git command failed: %w", err)
	}
	return nil
}

func gitCommandOutput(repoDir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git command failed: %w", err)
	}
	return output, nil
}
