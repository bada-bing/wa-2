package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"wa-2/arc"
	"wa-2/clockify"
	"wa-2/coldturkey"
	"wa-2/logseq"
	"wa-2/orb"
	"wa-2/sketchybar"
	"wa-2/util"
	"wa-2/vpn"
	"wa-2/vscode"
)

// this program will run inside of tmux session
// whenever you start a tmuxinator project it will be executed via `on_project_start` hook
// this is currently: ~/src/dotfiles/scripts/cui-next-work-session.sh
// that would essentially mean that all project would use the same `start_work_session script`

// 1. tn -> select `work_session` -> start tmuxinator project -> run on_project_start go application

// Notes
// useful hints: sudo visudo to add command to list of sudoers (e.g., kubectl)

// IMPR
// todo add github or gitlab url
// todo add jira url
// todo add linear url

// IMPR - open different localhost pages, depending on the scope
// arc tab create https://dev.wescale.io/shopping/cart/#/details/active
// arc tab create https://dev.wescale.io/shopping/#/carts/active

// IMPR - fetch URLs from Raindrop API (instead of having hard-coded versions in the source code)
// todo find localohost and DEV urls from raindrop and open them in browser
// todo read url from raindrop based on the session name

type Workflow_Config struct {
	Modules  []string `json:"modules"`
	Clockify struct {
		Project string `json:"project"`
	} `json:"clockify"`
}

func findConfigFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Define the relative path
	relativePath := "src/go/src/wa-2/data/cart-ui-next.wa-2.json"

	// Create the full path
	fullPath := filepath.Join(homeDir, relativePath)
	return fullPath
}

func readConfig(configFile string) Workflow_Config {
	// Read the JSON file
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON into the struct
	var config Workflow_Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// TODO validate using json schema

	return config
}

func main() {

	fullPath := findConfigFile()
	config := readConfig(fullPath)

	issueKey := util.GetIssueKey("")
	fmt.Printf("Active Issue: %s \n", issueKey)

	// Print the modules
	fmt.Println("Modules:", config.Modules)

	fmt.Println("📟 Start work session")

	var wg sync.WaitGroup
	wg.Add(8) // TODO make dynamic (based on number of modules)

	go orb.StartOrbStack(&wg)

	// // TODO don't start VPN if connected to corporate wifi
	go vpn.StartOpenVPNConnection(&wg) // tunnelblick
	go arc.SetupSpace(&wg)
	go clockify.StartTimer(&wg, config.Clockify.Project)
	go vscode.OpenWorkspace(&wg, "")
	go sketchybar.LoadWorkConfiguration(&wg)
	go logseq.OpenActivePage(issueKey)

	go coldturkey.StartWorkSession(&wg)

	wg.Wait()
}
