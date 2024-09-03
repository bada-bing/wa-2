package util

import (
	"log"
	"os/exec"
	"strings"
)

func GetIssueKey(dirname string) string {
	if dirname == "" {
		dirname = "~/src/cart-ui-next"
	}

	// TODO use os home dir and path join instead of hardcoded location
	// TODO move and rename issueKey.sh script
	cmd := exec.Command("bash", "-c", "/Users/miki/src/wa-2/clockify/issueKey.sh", dirname)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Error trying to get the issue key: ", err)
	}

	return strings.TrimSpace(string(output))
}
