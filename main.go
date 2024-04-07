package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type Licence struct {
	Name string `json:"name"`
}

type Repository struct {
	FullName    string  `json:"full_name"`
	Name        string  `json:"name"`
	Language    string  `json:"language"`
	Description string  `json:"description"`
	CloneUrl    string  `json:"clone_url"`
	Stars       int     `json:"stargazers_count"`
	Forks       int     `json:"forks"`
	Issues      int     `json:"open_issues"`
	Licence     Licence `json:"licence"`
}

type QueryResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []Repository `json:"items"`
}

func main() {
	args := os.Args
	isFork := flag.Bool("f", false, "Force search")
	flag.Parse()

	if *isFork && !isGitHubCLIInstalled() {
		fmt.Println("gh (github-cli) is not installed. Please install it before using the -f flag.")
		os.Exit(1)
	}

	queryIndex := findFlagIndex(args[1:]) + 1

	if len(args) < 2 {
		fmt.Println("Usage:\nCloning: rlu <search_query>\nForking: Cloning: rlu -f <search_query>")
		os.Exit(1)
	}

	client := http.Client{}
	url := "https://api.github.com/search/repositories?q=" + args[queryIndex] + "&sort=stars" + "&order=desc"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/vnd.github.text-match+json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Unexpected status: ", res.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		os.Exit(1)
	}

	var parsedResp QueryResponse
	if err := json.Unmarshal(body, &parsedResp); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	if parsedResp.TotalCount == 0 {
		fmt.Println("No repositories found")
		os.Exit(0)
	}

	selectedItem := DisplayTui(parsedResp.Items)

	if selectedItem.cloneUrl == "" {
		println("No repo selected")
		os.Exit(0)
	}

	{
		destination := getDestinationDir(selectedItem.shortTitle)

		command := buildCmdBase(*isFork)

		if *isFork {
			command = append(command, selectedItem.fullName)
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

	os.Exit(0)
}

func isGitHubCLIInstalled() bool {
	cmd := exec.Command("which", "gh")
	err := cmd.Run()
	return err == nil
}

func findFlagIndex(args []string) int {
	for i, arg := range args {
		if arg[0] != '-' {
			return i
		}
	}
	return len(args)
}

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
