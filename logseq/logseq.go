package logseq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// - Assume LogSeq is already running and http server is on
// - you want to open the page bootstrapped in wa-1

// 1. identify the logseq page title (based on the issue key)
// 2. send the pushState curl request (to logseq http server)

// 1. Identify the task (e.g., Jira Issue or Book Chapter)
// 1.1 - extract the ISSUE key from the current branch
// check GetIssueKey from clockify.go (extract it and make it util function)

// TODO consider to extract the following piece into a shell script, since it is used by both WA1 and WA2
// consider does not need that you should do, but think about it
// - TODO if LogSeq is not running, start it and start the HTTP server automatically

// # 1.2 - identify the logseq page title

type RequestBody struct {
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
}

func pushStateToLogseq(pageTitle string) int {
	logseqServer := "http://127.0.0.1:12315/api"
	logseqServerAPIToken := os.Getenv("LOGSEQ_SERVER_API_TOKEN")

	reqBody := RequestBody{
		Method: "logseq.app.pushState",
		Args:   []interface{}{"page", map[string]string{"name": pageTitle}},
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return 500
	}

	pushStateReq, err := http.NewRequest("POST", logseqServer, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating pushState request:", err)
		return 500
	}

	pushStateReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", logseqServerAPIToken))
	pushStateReq.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	pushStateResp, err := client.Do(pushStateReq)

	if err != nil {
		fmt.Println("Error sending pushState request:", err)
		return 500
	}
	defer pushStateResp.Body.Close()

	statusCode := pushStateResp.StatusCode

	return statusCode
}

func OpenActivePage(issueKey string) {

	if issueKey == "" {
		fmt.Println("issueKey is empty. return")
		return
	}

	// TODO you only make Jira Request for WPR tsks, for learning projects not ❗
	issue := getJiraIssue(issueKey)
	issueSummary := getIssueSummary(issue)

	pageTitle := fmt.Sprintf("%s-%s", issueKey, issueSummary)
	fmt.Println("Page Title:", pageTitle)

	// 2. - send the pushState curl request - to open the correct page in Logseq
	// Make the pushState request
	statusCode := pushStateToLogseq(pageTitle)

	fmt.Printf("Open Logseq Page: %s; Status: %d\n", pageTitle, statusCode)
}
