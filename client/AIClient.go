package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func AskAI(request OllamaRequest) (string, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", nil
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var msg ollamaResponse
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		fmt.Print(msg.Response)
		response += msg.Response

		if msg.Done {
			break
		}
	}

	return response, nil
}
