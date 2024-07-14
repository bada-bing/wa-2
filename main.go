package main

import (
	"sync"
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

// 0. Print "Start work session"
// 1. Identify the project (work or study), to figure out what tools and how to run them
// 2. Start necessary tools
// - 2.1 Start OrbStack ✅
// - 2.2 Start OpenVPN ✅
// - 2.3 Arc Browser Management
// - - - clear the space
// - - - open only relevant browsers
// - 2.4 Unit Tests, linting... ❓
// - 2.5. Open the correct VS Code repository
// - 2.6. Open the correct Logseq Page
// 3. Start Clockify Session

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	orb.StartOrbStack(&wg)
	vpn.StartOpenVPNConnection(&wg)
	wg.Wait()
}
