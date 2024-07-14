package arc

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// "\" in url breaks the code - at least it was breaking the arc cli

type Tab struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	WindowID int    `json:"windowId"`
	TabID    int    `json:"tabId"`
	Location string `json:"location"` // can be unpinned, pinned or topApp
}

// 1. call applescript
// 2. remove quotes from the output // there are no extra quotes in string in go program
// 3. unescape the output // json.Unmarshal does this automatically
// 4. convert the output into an array of tab structs
func getArcTabs() []Tab {

	cmd := exec.Command("osascript", "./arc/GetArcTabs.scpt")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	var tabs []Tab
	err = json.Unmarshal(output, &tabs)
	if err != nil {
		log.Fatalf("json.Unmarshal failed with %s\n", err)
	}

	return tabs
}

func findSimilarTabs(match string) []Tab {
	var tabs []Tab

	for _, tab := range getArcTabs() {
		if strings.Contains(tab.URL, match) {
			tabs = append(tabs, tab)
		}
	}

	return tabs
}

// todo: improve - only consider unpinned tab (and ignore any pinned or topApp tabs)
// todo: improve - use recursion as starting point as well, instead of passing initial set of tabs
// this function does not need to return
func clearTabs(match string) {
	tabs := findSimilarTabs(match)
	for len(tabs) > 0 {
		if len(tabs) == 1 {
			if !strings.Contains(match, "localhost") {
				return
			}

			// if only one is remaining, focus that tab and reload it
			cmd := exec.Command("arc", "tab", "select", strconv.Itoa(tabs[0].TabID))
			err := cmd.Run()
			if err != nil {
				log.Fatalf("[Select Tab] cmd.Run() failed with %s\n", err)
			}
			cmd = exec.Command("arc", "tab", "reload")
			err = cmd.Run()
			if err != nil {
				log.Fatalf("[Reload Tab] cmd.Run() failed with %s \n", err)
			}
			return
		}

		tab := tabs[0]
		fmt.Printf("Delete %s tab - %d\n", match, tab.TabID)
		cmd := exec.Command("arc", "tab", "close", strconv.Itoa(tab.TabID))
		err := cmd.Run()
		if err != nil {
			log.Fatalf("[Close Tab] cmd.Run() failed with %s [%d]\n", err, tab.TabID)
		}
		tabs = findSimilarTabs(match) // arc changes the ids of tabs after each tab is closed
	}

	// todo make it generic.. if no starting tabs, only create localhost tab
	if !strings.Contains(match, "localhost") {
		return
	}
	// if no tabs, create one
	cmd := exec.Command("arc", "tab", "create", "http://localhost:3000/shopping/cart/#/details/active")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[Create Tab] cmd.Run() failed with %s\n", err)
	}

}

func SetupSpace(wg *sync.WaitGroup) {
	defer wg.Done()

	// tabs := findSimilarTabs("localhost:3000")

	clearTabs("localhost:3000")
	clearTabs("https://dev.wescale.io/")
	fmt.Println("✔️ [arc] open development web pages")

	// Use work space
	// arc space focus 2 - no need for this; it automatically selects the correct space because of air traffic control

	// keep browser space tidy (use only one localhost and one dev page)
	// keep only one ✅ (identify all 'similar' tabs and only keep one or create a tab if none is open)
	// - kill all extra tabs ✅
	// - select and reload the remaining tab ✅
	// - if no tabs to start with, open new localhost fragment-dev-server ✅

	// impr: based on the scope of the story (from logseq) you can open the related page (e.g., overview or details)
	// impr: close all related unpinned jira links and open only one (or keep only one)
}
