package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your prompt: ")
	prompt, _ := reader.ReadString('\n')
	prompt = strings.TrimSpace(prompt)

	reqBody := OllamaRequest{
		Model:  "llama3",
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
