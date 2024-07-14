package util

import (
	"log"
	"os/exec"
	"strings"
)

func GetIssueKey() string {
	// TODO use os home dir and path join instead of hardcoded location
	cmd := exec.Command("bash", "-c", "/Users/miki/src/go/src/wa-2/clockify/issueKey.sh")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(output))
}
