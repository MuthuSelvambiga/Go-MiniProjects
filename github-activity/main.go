package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// defining struct for the fields we need
type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		RefType string     `json:"reftype"` //for createEvent
		Action  string     `json:"action"`  //for issues
		Commits []struct{} `json:"commits"` //for pushevents

	} `json:"payload"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}
	username := os.Args[1]
	fmt.Printf("Fetching recent activity for: %s\n", username)

	//1.Build the API URL
	url := "https://api.github.com/users/" + username + "/events"

	//2.Make HTTP GET Request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:Unable to fetch data from Github API")
		return
	}
	defer resp.Body.Close()

	//3.Check if request was successful
	if resp.StatusCode != 200 {
		fmt.Printf("Error:GitHub API returned status %d\n", resp.StatusCode)
		return
	}

	//4.Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:Unable to read response")
		return
	}

	//5.Print the raw JSON response
	//fmt.Println(string(body))

	//6.Parse json into []Event
	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	//7.Loop through and print human-readable messages
	for _, e := range events {
		switch e.Type {
		case "CreateEvent":
			refType := e.Payload.RefType
			if refType == "" {
				refType = "something"
			}

			if refType == "repository" {
				fmt.Printf("- Created a new repository %s\n", e.Repo.Name)
			} else {
				fmt.Printf("- Created %s in %s\n", refType, e.Repo.Name)
			}
		case "PushEvent":
			commitCount := len(e.Payload.Commits)
			fmt.Printf("- Pushed %d commits to %s\n", commitCount, e.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("- %s an issue in %s\n", e.Payload.Action, e.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf("- %s a pull request in %s\n", e.Payload.Action, e.Repo.Name)
		case "WatchEvent":
			fmt.Printf("- Starred %s\n", e.Repo.Name)
		default:
			fmt.Printf("- %s in %s\n", e.Type, e.Repo.Name)

		}
	}
}
