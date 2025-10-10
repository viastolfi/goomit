package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"main/prompt"
)

func GenerateConfig() (string, error) {
	var context string
	cmd := exec.Command("mkdir", ".goomit")
	_, err := cmd.CombinedOutput()

	// TODO: make this :
	// - platform agnostic
	// - less 'harcoded'
	// - working even if you are not at the root dir of the project
	if err != nil {
		if err.Error() == "exit status 1" {
			fmt.Println(".goomit/ dir already exist")
		} else {
			return "", fmt.Errorf("Error while adding file to commit \n", err)
		}
	} else {
		fmt.Println("config dir created")
	}

	cmd = exec.Command("ls", "README.md")
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("No README.md file detected")
	} else {
		fmt.Println("README.md file use in config generation")
		b, err := os.ReadFile("README.md")
		if err != nil {
			return "", err
		}
		context += string(b)
	}

	repo, err := getRepoName()
	if err != nil {
		return "", err
	}

	fmt.Println("Repo find :", repo, "trying to access it via github API")

	opts, err := githubApiCall(repo)
	if err != nil {
		return "", err
	}

	context += opts[0]
	context += opts[1]

	p := prompt.GenerateConfPrompt(context)

	return p, nil
}

func getRepoName() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error: could not read remote.origin.url — is this a Git repo?")
	}

	url := strings.TrimSpace(out.String())

	re := regexp.MustCompile(`(?i)[:/]([^/]+/[^/]+?)(?:\.git)?$`)
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", fmt.Errorf("Could not parse owner/repo from:", url)
	}

	repo := match[1]
	return repo, nil
}

func githubApiCall(repo string) ([]string, error) {
	var out []string

	url := "https://api.github.com/repos/" + repo

	fmt.Println("Performing ", url, "API call")
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error while getting github info using api :", err)
	}

	var j map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response body: ", err)
	}

	lang, _ := j["language"].(string)
	desc, _ := j["description"].(string)

	fmt.Println("Repo detected language:", lang)
	fmt.Println("Repo detected description:", desc)

	out = append(out, "language: "+lang, "description:"+desc)
	return out, nil
}

func WriteConfig(text string) error {
	f, err := os.Create(".goomit/context.md")
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(text))
	if err != nil {
		panic(err)
	}

	return nil
}
