package main

import (
	"fmt"
	"os"
	"os/exec"
)

func getDestinationDir(dirName string) string {

	cwd, err := os.Getwd()

	if err != nil {
		println("Error getting current working directory")
	}
	return cwd + "/" + dirName
}

func buildCmdBase(isFork bool) []string {
	if isFork {
		return []string{"gh", "repo", "fork"}
	}
	return []string{"git", "clone"}

}

func CloneRepo(selectedItem item, fork bool) {
	if selectedItem.cloneUrl == "" {
		println("No repo selected")
		os.Exit(0)
	}

	destination := getDestinationDir(selectedItem.shortTitle)

	command := buildCmdBase(fork)

	if fork {
		command = append(command, selectedItem.fullName, destination)
	} else {
		command = append(command, selectedItem.cloneUrl, destination)
	}

	cmd := exec.Command(command[0], command[1:]...)

	// Connect to system IO
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		os.Exit(1)
	}
}
