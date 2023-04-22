package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	filePath := "REAMD.md"
	commitMsg := "Auto commit from Go"

	// run git add instruction
	gitAddCmd := exec.Command("git", "add", filePath)
	gitAddCmd.Stdout = os.Stdout
	gitAddCmd.Stderr = os.Stderr

	if err := gitAddCmd.Run(); err != nil {
		fmt.Println("Error: git commit failed")
		os.Exit(1)
	}

	// run git commit instruction
	gitCommitCmd := exec.Command("git", "commit", "-m", commitMsg)
	gitCommitCmd.Stdout = os.Stdout
	gitCommitCmd.Stdin = os.Stderr
	if err := gitCommitCmd.Run(); err != nil {
		fmt.Println("Error: git commit failed")
		os.Exit(1)
	}

	// get the presenting branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchCmdOut, err := branchCmd.Output()
	if err != nil {
		fmt.Println("Error: failed to get current branch")
		os.Exit(1)
	}
	branch := strings.TrimSpace(string(branchCmdOut))

	gitPushCmd := exec.Command("git", "push", "origin", branch)
	gitPushCmd.Stdout = os.Stdout
	gitPushCmd.Stderr = os.Stderr

	if err := gitPushCmd.Run(); err != nil {
		fmt.Println("Error: git push failed")
		os.Exit(1)
	}

	fmt.Println("Successfully pushed changes to Github")

}
