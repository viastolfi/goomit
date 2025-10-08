package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"main/client"
	"main/config"
	"main/prompt"

	"golang.org/x/term"
)

func getGitDiff() (string, error) {
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

func commit(message string) error {
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

func main() {
	var modelName string = "llama3"

	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "-m" || args[0] == "--model" {
			modelName = args[1]
		} else if args[0] == "config" && args[1] == "generate" {
			if err := config.GenerateConfig(); err != nil {
				log.Fatalf("Error while generating config\n%s", err)
			}
			fmt.Println("Config generated on '.goomit/'")
			return
		} else {
			log.Fatalf("unknown argument : %s \nPlease refer to `goomit --help` for more help", args[0])
		}
	}

	diff, err := getGitDiff()
	if err != nil {
		log.Fatalf("Error while getting git diff \n %s", err)
	}

	prompt := prompt.GeneratePrompt(diff)
	fmt.Println(prompt)

	fmt.Println("GENERATING...")

	reqBody := client.OllamaRequest{
		Model:  modelName,
		Prompt: prompt,
	}

	resp, err := client.AskAI(reqBody)
	if err != nil {
		log.Fatalf("Error while generating", err)
	}

	fmt.Println("\nDo you want to use this commit message ? y/n/Y/N [y] : ")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	term.Restore(int(os.Stdin.Fd()), oldState)
	if b[0] == 'y' || b[0] == 'Y' || b[0] == '\r' {
		if err := commit(resp); err != nil {
			log.Fatalf("Error during commit phase : %s\nTry to commit yourself", err)
		}
		fmt.Println("All done ! You can now push your commits\nThanks for using goomit")
	}

	fmt.Println("Closing !")
}
