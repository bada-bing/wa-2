package vscode

import (
	"fmt"
	"log"
	"os/exec"
)

func getCodeWorkspacePath(projectName string) string {
	dropboxDir := os.Getenv("DROPBOX_DIR")
	// basePath should be configurable (or at least read from ENV)
	basePath := dropboxDir + "/toolbox/tools/vs_code/workspaces/"
	return basePath + projectName + ".code-workspace"
}

func OpenWorkspace(wg *sync.WaitGroup, projectName string) {
	defer wg.Done()

	if projectName == "" {
		projectName = "cart-ui-next"
	}

	codeWorkspace := getCodeWorkspacePath(projectName)
	openWorkspace := "code " + codeWorkspace

	var cmd = exec.Command("bash", "-c", openWorkspace)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start VS Code Workspace: %v", err)
	} else {
		fmt.Println("✔️ [vscode] open workspace ", projectName)
	}

}



