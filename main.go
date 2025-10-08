package main

import (
	"fmt"
	"log"
	"os"

	"main/client"
	"main/config"
	"main/prompt"

	"golang.org/x/term"
)

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

	diff, err := client.GetGitDiff()
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
		if err := client.Commit(resp); err != nil {
			log.Fatalf("Error during commit phase : %s\nTry to commit yourself", err)
		}
		fmt.Println("All done ! You can now push your commits\nThanks for using goomit")
	}

	fmt.Println("Closing !")
}
