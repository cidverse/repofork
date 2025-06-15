package fork

import (
	"fmt"
	"os"
	"os/exec"
)

type Git struct {
	repoDir string
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
