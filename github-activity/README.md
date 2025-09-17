# GitHub Activity CLI 
A simple command line interface (CLI) to fetch the recent activity of a GitHub user and display it in the terminal.

## GOAL
- Accept a GitHub username as an argument
- Fetch recent activity using the GitHub API: `https://api.github.com/users/<username>/events`
- Display activity in a human-readable format
- Handle errors gracefully (invalid usernames, no activity, API failures)
- **No external libraries** — only Go standard library

## Usage

github-activity <username>
