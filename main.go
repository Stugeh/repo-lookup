package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	if len(args) < 2 {
		fmt.Println("Usage: rlu <search_query>")
		os.Exit(1)
	}

	client := http.Client{}
	url := "https://api.github.com/search/repositories?q=" + args[1] + "&sort=stars" + "&order=desc"

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

	DisplayTui(parsedResp.Items)
}
