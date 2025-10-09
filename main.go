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
	var request client.OllamaRequest
	request.Model = "llama3"

	args := os.Args[1:]
	if len(args) > 0 {
		for i := range args {
			if args[i] == "-m" || args[i] == "--model" {
				if i == len(args)-1 {
					log.Fatalf("Wrong inline argument usage\nuse --help for more information")
				}
				request.Model = args[i+1]
				p, err := prompt.GeneratePrompt()
				if err != nil {
					log.Fatalf("Error while generating prompt\n", err)
				}
				request.Prompt = p
			}
			if args[i] == "config" {
				if i == len(args)-1 {
					log.Fatalf("Wrong inline argument usage\nuse --help for more information")
				} else if args[i+1] == "generate" {
					p, err := config.GenerateConfig()
					if err != nil {
						log.Fatalf("Error while generating config prompt\n", err)
					}
					request.Prompt = p
				}
			}
		}
	} else {
		p, err := prompt.GeneratePrompt()
		if err != nil {
			log.Fatalf("Error while generating config prompt\n", err)
		}
		request.Prompt = p
	}

	fmt.Println(request.Prompt)

	fmt.Println("GENERATING...")

	resp, err := client.AskAI(request)
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
