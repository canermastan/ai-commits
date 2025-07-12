package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/canermastan/ai-commits/config"
	"github.com/canermastan/ai-commits/internal/ai"
	"github.com/canermastan/ai-commits/internal/git"
	"github.com/canermastan/ai-commits/internal/ui"
)

type options struct {
	fastMode bool
}

func parseFlags() options {
	opts := options{}
	flag.BoolVar(&opts.fastMode, "fast", false, "Skip user interaction and generate commit message directly from diff")
	flag.Parse()
	return opts
}

func handleUnstagedChanges() {
	unstagedFiles, err := git.GetUnstagedFiles()
	if err != nil {
		ui.ShowError("getting unstaged files: %v", err)
	}

	fmt.Println("\nUnstaged files:")
	ui.ShowFiles(unstagedFiles)
	fmt.Println("\nPlease add the files to the staging area and try again.")
	os.Exit(1)
}

func validateEnvironment() {
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		ui.ShowError("GEMINI_API_KEY environment variable is not set. Please set it and try again.")
		os.Exit(1)
	}
}

func getStagedFiles() []string {
	files, err := git.GetStagedFiles()
	if err != nil {
		handleUnstagedChanges()
	}
	return files
}

func getExplanation(fastMode bool) string {
	if fastMode {
		return "Generate a commit message based on the changes in the diff"
	}

	explanation, err := ui.GetExplanation()
	if err != nil {
		ui.ShowError("reading explanation: %v", err)
	}
	return explanation
}

func getDiff() string {
	diff, err := git.GetDiff()
	if err != nil {
		ui.ShowError("getting diff: %v", err)
	}

	maxLen := config.GetMaxDiffSize()
	if len(diff) > maxLen {
		fmt.Println("Warning: Diff too large, using partial diff.")
		return diff[:maxLen]
	}
	return diff
}

func generateCommitMessage(explanation, diff string, fastMode bool) string {
	prompt := ai.BuildPrompt(explanation, diff)

	var message string
	var err error

	if fastMode {
		message, err = ui.WithFastLoading(func() (string, error) {
			return ai.CallAI(prompt)
		})
	} else {
		message, err = ui.WithLoading(func() (string, error) {
			return ai.CallAI(prompt)
		})
	}

	if err != nil {
		ui.ShowError("calling AI: %v", err)
	}
	return message
}

func handleCommit(message string) {
	if err := git.Commit(message); err != nil {
		ui.ShowError("creating commit: %v", err)
	}
	ui.ShowSuccess("Commit created successfully!")
}

func Root() {
	opts := parseFlags()
	validateEnvironment()

	files := getStagedFiles()
	if !opts.fastMode {
		fmt.Println("\nStaged files:")
		ui.ShowFiles(files)
	}

	explanation := getExplanation(opts.fastMode)
	diff := getDiff()
	message := generateCommitMessage(explanation, diff, opts.fastMode)

	if opts.fastMode {
		fmt.Println(message)
		return
	}

	confirm, err := ui.ConfirmCommit(message)
	if err != nil {
		ui.ShowError("reading confirmation: %v", err)
	}

	if !confirm {
		fmt.Println("Commit cancelled.")
		return
	}

	handleCommit(message)
}
