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
