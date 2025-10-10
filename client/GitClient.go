package client

import (
	"fmt"
	"os/exec"
)

func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--color=always")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("No git repo detected in your current directory \n", err)
	}

	if len(output) == 0 {
		return "", fmt.Errorf("No changes detected, are you sure you have something to commit ?")
	}

	return string(output), nil
}

func Commit(message string) error {
	cmd := exec.Command("git", "add", "-u")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error while adding file to commit \n", err)
	}

	fmt.Println("File added!")
	cmd = exec.Command("git", "commit", "-m", message)
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Errorf("Error while commiting changes \n", err)
	}

	fmt.Println("Commit created!")
	return nil
}
