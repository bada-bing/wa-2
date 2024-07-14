package logseq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Issue struct {
	Fields Fields `json:"fields"`
}

type Fields struct {
	Summary string `json:"summary"`
}

func createJiraRequest(issueKey string) *http.Request {
	jiraAccessToken := os.Getenv("JIRA_ACCESS_TOKEN")

	url := fmt.Sprintf("https://wescalehq.atlassian.net/rest/api/3/issue/%s", issueKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating network request")
		return &http.Request{}
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("accept-language", "en-US,en")
	req.Header.Add("authorization", fmt.Sprintf("Basic %s", jiraAccessToken))
	// req.Header.Add("user-agent", "vscode-restclient")

	return req
}

func getJiraIssue(issueKey string) Issue {
	req := createJiraRequest(issueKey)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making network request", err)
		return Issue{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return Issue{}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return Issue{}
	}

	var issue Issue
	err = json.Unmarshal(body, &issue)

	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return Issue{}
	}

	return issue
}

func getIssueSummary(issue Issue) string {
	fmt.Println("Task Summary:", issue.Fields.Summary)

	// Process the task summary string
	issueSummary := strings.ToLower(issue.Fields.Summary)
	issueSummary = strings.ReplaceAll(issueSummary, " ", "-")
	issueSummary = strings.ReplaceAll(issueSummary, "/", "_")
	issueSummary = strings.ReplaceAll(issueSummary, "shopping-cart", "cart")
	issueSummary = strings.ReplaceAll(issueSummary, "-the-", "-")
	issueSummary = strings.ReplaceAll(issueSummary, "\"", "")

	return issueSummary
}
