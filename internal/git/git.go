package git

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

var (
	ErrNoStagedFiles      = errors.New("no staged files")
	ErrNoUnstagedFiles    = errors.New("no unstaged files")
	ErrEmptyCommitMessage = errors.New("commit message cannot be empty")
)

// GetStagedFiles returns a list of files that are staged for commit.
// Returns ErrNoStagedFiles if there are no staged files.
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil, ErrNoStagedFiles
	}

	files := strings.Split(trimmed, "\n")
	return files, nil
}

// GetUnstagedFiles returns a list of files that have changes but are not staged.
// Returns ErrNoUnstagedFiles if there are no unstaged changes.
func GetUnstagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil, ErrNoUnstagedFiles
	}

	files := strings.Split(trimmed, "\n")
	return files, nil
}

// GetDiff returns the diff of all staged changes.
// The diff includes only the staged files (--cached).
func GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

// Commit creates a new commit with the given message.
func Commit(message string) error {
	if message == "" {
		return ErrEmptyCommitMessage
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
