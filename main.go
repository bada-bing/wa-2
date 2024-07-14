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
	"wa-2/orb"
	"wa-2/vpn"
)

// this small program will be run inside of tmux session
// whenever you start a tmuxinator project it will be executed via `on_project_start` hook
// this is currently: ~/src/dotfiles/scripts/cui-next-work-session.sh
// that would essentially mean that all project would use the same `start_work_session script`

// 1. tn -> select `work_session` -> start tmuxinator project -> run on_project_start go application

// go application will determine what modules to start (e.g., don't start vpn if it is not a work session)
// even turn off vpn if it is a study session
// this should be defined via configuration
// { arc: {localhost: "localhost:3000/shopping/cart/#/"}, vpn: on, orb: on, logseq: ...}

// 0. Print "Start work session"
// 1. Identify the project (work or study), to figure out what tools and how to run them
// this could be sent as a parameter (from tmuxinator project)
// 1.1. Define Configuration for different projects
// 2. Start necessary tools
// - 2.1 Start OrbStack ✅
// - 2.2 Start OpenVPN ✅
// - 2.3 Arc Browser Management ✅
// - - - clear the space
// - - - - for starters: delete all localhost pages
// - - - open only relevant pages
// - 2.4. Open the correct VS Code repository
// - - 2.4.1 open the workspace based on the related project
// - - 2.4.❓ keep only a couple of editor panes (e.g., last 2, 3) // not fully possible
// - - - don't bother too much with this, because you are already in the realm of the work session (since you are close to the source code once in workspace)
// - - - what you can do is something like open empty workspace, and then `code <file> --reuse-window`
// - 2.5 Open the correct Logseq Page
// - 2.6. Unit Tests, linting... ❓
// 3. Start Clockify Session ✅
// 4. Focus
// - 4.1 cold turkey work_session block
// - 4.2 dnd focus profile macos
// - 4.❓ mute slack and teams pages (close teams app, ...)
// - - arc doesn't support mute, I would have to do that using javascript (it is probably not that straightforward)
// - - ? do you need slack as favorite or pinned tab at all
// - - you could create a protocol to have it closed, and open it a couple of times a day or when you receive a phone application
// - - run raycast script which opens unreads directly so you can process them

type Workflow_Config struct {
	Modules []string `json:"modules"`
}

func findConfigFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	fmt.Println("Home directory:", homeDir)

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

	// Print the modules
	fmt.Println("Modules:", config.Modules)

	fmt.Println("📟 Start work session")

	var wg sync.WaitGroup
	wg.Add(4)

	go orb.StartOrbStack(&wg)
	go vpn.StartOpenVPNConnection(&wg)
	go arc.SetupSpace(&wg)
	go clockify.StartTimer(&wg)

	wg.Wait()
}
