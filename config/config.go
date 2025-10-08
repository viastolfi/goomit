package config

import (
	"fmt"
	"os/exec"
)

func GenerateConfig() error {
	var files []string

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
			return fmt.Errorf("Error while adding file to commit \n", err)
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
		files = append(files, "README.md")
	}

	return nil
}
