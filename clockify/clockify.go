package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const WORKSPACES_URL = "https://api.clockify.me/api/v1/workspaces"

var MAIN_WORKSPACE_ID = os.Getenv("CLOCKIFY_MAIN_WORKSPACE_ID")

var STORIES_ID = os.Getenv("STORIES_ID") // jira project

type Issue struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

// currently only works for cart-ui-next
func GetIssueKey() string {
	// TODO use os home dir and path join instead of hardcoded location
	cmd := exec.Command("bash", "-c", "/Users/miki/src/go/src/wa-2/clockify/issueKey.sh")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(output))
}

func getCurrentJiraIssue() Issue {

	tasksUrl := fmt.Sprintf("%s/%s/projects/%s/tasks", WORKSPACES_URL, MAIN_WORKSPACE_ID, STORIES_ID)

	apiKey := os.Getenv("CLOCKIFY_API_KEY")

	client := &http.Client{}
	req, err := http.NewRequest("GET", tasksUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var clockifyIssues []Issue
	err = json.Unmarshal(body, &clockifyIssues)
	if err != nil {
		log.Fatal(err)
	}

	currentIssueKey := GetIssueKey()

	for _, issue := range clockifyIssues {
		// fmt.Printf("name: %s, id: %s\n", task.Name, task.Id)
		if issue.Name == currentIssueKey {
			return issue
		}
	}
	return Issue{
		Id:   currentIssueKey,
		Name: "doesnt_exist",
	}
}

func getFormattedTime() string {
	currentTime := time.Now()
	twoHoursAgo := currentTime.Add(-2 * time.Hour)
	formattedTime := twoHoursAgo.Format("2006-01-02T15:04:05Z")
	return formattedTime
}

func StartTimer(wg *sync.WaitGroup) {
	defer wg.Done()

	apiKey := os.Getenv("CLOCKIFY_API_KEY")

	formattedTime := getFormattedTime()
	issue := getCurrentJiraIssue()

	if issue.Name == "doesnt_exist" {
		fmt.Printf("❌ [clockify] %s issue does not exist. Did you bootstap it properly?\n", issue.Id)
		return

		// TODO you don't need to exit here, just don't use task id in the request body down below
	}

	description := "✨ Programming"

	ADD_ENTRY_URL := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/time-entries", MAIN_WORKSPACE_ID)

	type RequestBody struct {
		Billable    bool   `json:"billable"`
		Description string `json:"description"`
		ProjectID   string `json:"projectId"`
		TaskID      string `json:"taskId"`
		Start       string `json:"start"`
		Type        string `json:"type"`
	}

	body := &RequestBody{
		Billable:    true,
		Description: description,
		ProjectID:   STORIES_ID,
		TaskID:      issue.Id,
		Start:       formattedTime,
		Type:        "REGULAR",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", ADD_ENTRY_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	fmt.Printf("Add time entry: Issue: %s; Status: %d\n", issue.Name, statusCode)
}
