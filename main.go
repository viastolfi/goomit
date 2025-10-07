package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

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

func main() {
	var modelName string = "llama3"

	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "-m" || args[0] == "--model" {
			modelName = args[1]
		} else {
			log.Fatalf("unknown argument : %s \nPlease refer to `goomit --help` for more help", args[0])
		}
	}

	diff, err := getGitDiff()
	if err != nil {
		log.Fatalf("Error while getting git diff \n %s", err)
	}

	fmt.Println(diff)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your prompt: ")
	prompt, _ := reader.ReadString('\n')
	prompt = strings.TrimSpace(prompt)

	reqBody := OllamaRequest{
		Model:  modelName,
		Prompt: prompt,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var msg OllamaResponse
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		fmt.Print(msg.Response)

		if msg.Done {
			break
		}
	}
	fmt.Println("\n--- Generation Complete ---")
}
